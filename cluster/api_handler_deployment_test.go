package cluster

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/fake"
	"github.com/scalezilla/scalezilla/cri"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCluster_api_handler_deployment(t *testing.T) {
	assert := assert.New(t)

	t.Run("apply_basic", func(t *testing.T) {
		tests := []struct {
			method                                         string
			uri                                            string
			expectedStatusCode                             int
			expectedBody                                   string
			file, newFile                                  string
			createContainerError                           bool
			raftIsLeader, existsError                      bool
			decodeError, submitCommandDeploymentWriteError bool
		}{
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "",
				expectedStatusCode: 400,
			},
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "basic_success.hcl",
				expectedStatusCode: 403,
			},
			{
				method:                            "POST",
				uri:                               "/api/v1/deployment/apply",
				file:                              "basic_success.hcl",
				expectedStatusCode:                500,
				raftIsLeader:                      true,
				submitCommandDeploymentWriteError: true,
			},
			{
				method:               "POST",
				uri:                  "/api/v1/deployment/apply",
				file:                 "error_container_image.hcl",
				expectedStatusCode:   400,
				raftIsLeader:         true,
				createContainerError: true,
			},
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "error_bad_deployment_name_rfc.hcl",
				expectedStatusCode: 400,
				raftIsLeader:       true,
			},
			{
				method:               "POST",
				uri:                  "/api/v1/deployment/apply",
				file:                 "basic_success.hcl",
				expectedStatusCode:   400,
				raftIsLeader:         true,
				createContainerError: true,
			},
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "basic_success.hcl",
				newFile:            "basic_redis_success.hcl",
				expectedStatusCode: 200,
				raftIsLeader:       true,
			},
		}

		header := map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}

		for _, tc := range tests {
			cfg := basicClusterConfig{randomPort: false, dev: true}
			cluster := makeBasicCluster(cfg)
			cluster.config.DataDir = filepath.Join(cluster.config.DataDir, fake.CharactersN(10))
			store, err := cluster.buildStore()
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()
			assert.Nil(err)
			cluster.fsm = newFSM(store)

			var (
				payload  string
				dataFile []byte
			)
			if tc.file != "" {
				workingDir, err := os.Getwd()
				assert.Nil(err)

				configDir := "testdata/deployments"
				configFile := filepath.Join(workingDir, configDir, tc.file)
				dataFile, err = os.ReadFile(configFile)
				require.NoError(t, err)

				b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(dataFile)})
				require.NoError(t, err)
				payload = string(b)
			}

			mock := mockRafty{}
			if tc.raftIsLeader {
				mock.isLeader = true
			}
			cluster.rafty = &mock

			if tc.createContainerError {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return errors.New("container error")
				}
			} else {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return nil
				}
			}

			if tc.submitCommandDeploymentWriteError {
				cluster.di.deploymentEncodeCommandFunc = func(cmd deploymentState, w io.Writer) error {
					return errors.New("submit error")
				}
			}

			router := cluster.newApiRouters()
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, header, payload)
			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
		}
	})

	t.Run("apply_already_exists_get_error", func(t *testing.T) {
		tests := []struct {
			method                    string
			uri                       string
			expectedStatusCode        int
			expectedBody              string
			file, newFile             string
			createContainerError      bool
			raftIsLeader, existsError bool
			decodeError               bool
		}{
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "basic_success.hcl",
				newFile:            "basic_redis_success.hcl",
				expectedStatusCode: 500,
				raftIsLeader:       true,
			},
		}

		header := map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}

		for _, tc := range tests {
			cfg := basicClusterConfig{randomPort: false, dev: true}
			cluster := makeBasicCluster(cfg)
			cluster.config.DataDir = filepath.Join(cluster.config.DataDir, fake.CharactersN(10))
			store, err := cluster.buildStore()
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()
			assert.Nil(err)
			cluster.fsm = newFSM(store)

			var (
				payload  string
				dataFile []byte
			)
			if tc.file != "" {
				workingDir, err := os.Getwd()
				assert.Nil(err)

				configDir := "testdata/deployments"
				configFile := filepath.Join(workingDir, configDir, tc.file)
				dataFile, err = os.ReadFile(configFile)
				require.NoError(t, err)

				b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(dataFile)})
				require.NoError(t, err)
				payload = string(b)
			}

			mock := mockRafty{}
			if tc.raftIsLeader {
				mock.isLeader = true
			}
			cluster.rafty = &mock

			if tc.createContainerError {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return errors.New("container error")
				}
			} else {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return nil
				}
			}

			router := cluster.newApiRouters()
			cluster.fsm.memoryDeploymentExistsFunc = func(namespace, deploymentName []byte) bool {
				return true
			}

			cluster.fsm.memoryDeploymentGetFunc = func(namespace, deploymentName []byte) ([]byte, error) {
				return nil, errors.New("get deployment error")
			}

			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, header, payload)
			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
		}
	})

	t.Run("apply_already_exists_decode_error", func(t *testing.T) {
		tests := []struct {
			method               string
			uri                  string
			expectedStatusCode   int
			expectedBody         string
			file, newFile        string
			createContainerError bool
		}{
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "basic_success.hcl",
				newFile:            "basic_redis_success.hcl",
				expectedStatusCode: 500,
			},
		}

		header := map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}

		for _, tc := range tests {
			cfg := basicClusterConfig{randomPort: false, dev: true}
			cluster := makeBasicCluster(cfg)
			cluster.config.DataDir = filepath.Join(cluster.config.DataDir, fake.CharactersN(10))
			store, err := cluster.buildStore()
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()
			assert.Nil(err)
			cluster.fsm = newFSM(store)

			var (
				payload  string
				dataFile []byte
			)
			if tc.file != "" {
				workingDir, err := os.Getwd()
				assert.Nil(err)

				configDir := "testdata/deployments"
				configFile := filepath.Join(workingDir, configDir, tc.file)
				dataFile, err = os.ReadFile(configFile)
				require.NoError(t, err)

				b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(dataFile)})
				require.NoError(t, err)
				payload = string(b)
			}

			mock := mockRafty{}
			mock.isLeader = true
			cluster.rafty = &mock

			if tc.createContainerError {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return errors.New("container error")
				}
			} else {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return nil
				}
			}

			cluster.fsm.memoryDeploymentExistsFunc = func(namespace, deploymentName []byte) bool {
				return true
			}

			cluster.fsm.memoryDeploymentGetFunc = func(namespace, deploymentName []byte) ([]byte, error) {
				return []byte(`{a`), nil
			}

			router := cluster.newApiRouters()
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, header, payload)
			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
		}
	})

	t.Run("apply_already_exists_same_version", func(t *testing.T) {
		tests := []struct {
			method               string
			uri                  string
			expectedStatusCode   int
			expectedBody         string
			file, newFile        string
			createContainerError bool
		}{
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "basic_success.hcl",
				newFile:            "basic_redis_success.hcl",
				expectedStatusCode: 200,
			},
		}

		header := map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}

		for _, tc := range tests {
			cfg := basicClusterConfig{randomPort: false, dev: true}
			cluster := makeBasicCluster(cfg)
			cluster.config.DataDir = filepath.Join(cluster.config.DataDir, fake.CharactersN(10))
			store, err := cluster.buildStore()
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()
			assert.Nil(err)
			cluster.fsm = newFSM(store)

			var (
				payload  string
				dataFile []byte
			)
			if tc.file != "" {
				workingDir, err := os.Getwd()
				assert.Nil(err)

				configDir := "testdata/deployments"
				configFile := filepath.Join(workingDir, configDir, tc.file)
				dataFile, err = os.ReadFile(configFile)
				require.NoError(t, err)

				b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(dataFile)})
				require.NoError(t, err)
				payload = string(b)
			}

			mock := mockRafty{}
			mock.isLeader = true
			cluster.rafty = &mock

			if tc.createContainerError {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return errors.New("container error")
				}
			} else {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return nil
				}
			}

			cluster.fsm.memoryDeploymentExistsFunc = func(namespace, deploymentName []byte) bool {
				return true
			}

			cluster.fsm.memoryDeploymentGetFunc = func(namespace, deploymentName []byte) ([]byte, error) {
				spec, err := cluster.parseDeployment(dataFile)
				assert.NoError(err)
				replicaSetID := string(strings.ToLower(rand.Text())[:10])
				dc := deploymentContent{
					RawContent:   string(dataFile),
					Version:      1,
					CreatedAt:    time.Now(),
					ReplicaSetID: replicaSetID,
				}
				content := make(map[uint64]deploymentContent, 1)
				content[dc.Version] = dc
				state := deploymentState{
					Kind:               deploymentCommandSet,
					Name:               spec.Deployment.Name,
					NewRollingVersion:  int64(dc.Version),
					CurrentUsedVersion: dc.Version,
					Content:            content,
					MustBeStarted:      true,
				}

				cluster.fsm.memoryStore.deployment[spec.Deployment.Name] = state
				buffer := new(bytes.Buffer)
				err = deploymentEncodeCommand(state, buffer)
				assert.Nil(err)
				return buffer.Bytes(), nil
			}

			router := cluster.newApiRouters()
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, header, payload)
			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
		}
	})

	t.Run("apply_already_exists_new_version", func(t *testing.T) {
		tests := []struct {
			method               string
			uri                  string
			expectedStatusCode   int
			expectedBody         string
			file, newFile        string
			createContainerError bool
		}{
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "basic_success.hcl",
				newFile:            "basic_redis_success.hcl",
				expectedStatusCode: 200,
			},
		}

		header := map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}

		for _, tc := range tests {
			cfg := basicClusterConfig{randomPort: false, dev: true}
			cluster := makeBasicCluster(cfg)
			cluster.config.DataDir = filepath.Join(cluster.config.DataDir, fake.CharactersN(10))
			store, err := cluster.buildStore()
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()
			assert.Nil(err)
			cluster.fsm = newFSM(store)

			var (
				payload  string
				dataFile []byte
			)
			if tc.file != "" {
				workingDir, err := os.Getwd()
				assert.Nil(err)

				configDir := "testdata/deployments"
				configFile := filepath.Join(workingDir, configDir, tc.file)
				dataFile, err = os.ReadFile(configFile)
				require.NoError(t, err)

				b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(dataFile)})
				require.NoError(t, err)
				payload = string(b)
			}

			mock := mockRafty{}
			mock.isLeader = true
			cluster.rafty = &mock

			if tc.createContainerError {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return errors.New("container error")
				}
			} else {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return nil
				}
			}

			cluster.fsm.memoryDeploymentExistsFunc = func(namespace, deploymentName []byte) bool {
				return true
			}

			cluster.fsm.memoryDeploymentGetFunc = func(namespace, deploymentName []byte) ([]byte, error) {
				spec, err := cluster.parseDeployment(dataFile)
				assert.NoError(err)
				replicaSetID := string(strings.ToLower(rand.Text())[:10])
				dc := deploymentContent{
					RawContent:   string(dataFile),
					Version:      1,
					CreatedAt:    time.Now(),
					ReplicaSetID: replicaSetID,
				}
				content := make(map[uint64]deploymentContent, 1)
				content[dc.Version] = dc
				state := deploymentState{
					Kind:               deploymentCommandSet,
					Name:               spec.Deployment.Name,
					NewRollingVersion:  int64(dc.Version),
					CurrentUsedVersion: dc.Version,
					Content:            content,
					MustBeStarted:      true,
				}

				cluster.fsm.memoryStore.deployment[spec.Deployment.Name] = state
				buffer := new(bytes.Buffer)
				err = deploymentEncodeCommand(state, buffer)
				assert.Nil(err)
				return buffer.Bytes(), nil
			}

			workingDir, err := os.Getwd()
			assert.Nil(err)

			configDir := "testdata/deployments"
			configFile := filepath.Join(workingDir, configDir, tc.newFile)
			data, err := os.ReadFile(configFile)
			require.NoError(t, err)

			b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(data)})
			require.NoError(t, err)
			payload = string(b)

			router := cluster.newApiRouters()
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, header, payload)
			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
		}
	})

	t.Run("apply_already_exists_new_version_submit_error", func(t *testing.T) {
		tests := []struct {
			method               string
			uri                  string
			expectedStatusCode   int
			expectedBody         string
			file, newFile        string
			createContainerError bool
		}{
			{
				method:             "POST",
				uri:                "/api/v1/deployment/apply",
				file:               "basic_success.hcl",
				newFile:            "basic_redis_success.hcl",
				expectedStatusCode: 500,
			},
		}

		header := map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}

		for _, tc := range tests {
			cfg := basicClusterConfig{randomPort: false, dev: true}
			cluster := makeBasicCluster(cfg)
			cluster.config.DataDir = filepath.Join(cluster.config.DataDir, fake.CharactersN(10))
			store, err := cluster.buildStore()
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()
			assert.Nil(err)
			cluster.fsm = newFSM(store)

			var (
				payload  string
				dataFile []byte
			)
			if tc.file != "" {
				workingDir, err := os.Getwd()
				assert.Nil(err)

				configDir := "testdata/deployments"
				configFile := filepath.Join(workingDir, configDir, tc.file)
				dataFile, err = os.ReadFile(configFile)
				require.NoError(t, err)

				b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(dataFile)})
				require.NoError(t, err)
				payload = string(b)
			}

			mock := mockRafty{}
			mock.isLeader = true
			cluster.rafty = &mock

			if tc.createContainerError {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return errors.New("container error")
				}
			} else {
				cluster.di.createContainerFunc = func(ctx context.Context, spec cri.CreateContainerSpec) error {
					return nil
				}
			}

			cluster.fsm.memoryDeploymentExistsFunc = func(namespace, deploymentName []byte) bool {
				return true
			}

			cluster.fsm.memoryDeploymentGetFunc = func(namespace, deploymentName []byte) ([]byte, error) {
				spec, err := cluster.parseDeployment(dataFile)
				assert.NoError(err)
				replicaSetID := string(strings.ToLower(rand.Text())[:10])
				dc := deploymentContent{
					RawContent:   string(dataFile),
					Version:      1,
					CreatedAt:    time.Now(),
					ReplicaSetID: replicaSetID,
				}
				content := make(map[uint64]deploymentContent, 1)
				content[dc.Version] = dc
				state := deploymentState{
					Kind:               deploymentCommandSet,
					Name:               spec.Deployment.Name,
					NewRollingVersion:  int64(dc.Version),
					CurrentUsedVersion: dc.Version,
					Content:            content,
					MustBeStarted:      true,
				}

				cluster.fsm.memoryStore.deployment[spec.Deployment.Name] = state
				buffer := new(bytes.Buffer)
				err = deploymentEncodeCommand(state, buffer)
				assert.Nil(err)
				return buffer.Bytes(), nil
			}

			cluster.di.deploymentEncodeCommandFunc = func(cmd deploymentState, w io.Writer) error {
				return errors.New("submit error")
			}

			workingDir, err := os.Getwd()
			assert.Nil(err)

			configDir := "testdata/deployments"
			configFile := filepath.Join(workingDir, configDir, tc.newFile)
			data, err := os.ReadFile(configFile)
			require.NoError(t, err)

			b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(data)})
			require.NoError(t, err)
			payload = string(b)

			router := cluster.newApiRouters()
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, header, payload)
			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
		}
	})
}

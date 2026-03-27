package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/scalezilla/scalezilla/cri"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCluster_api_handler_deployment(t *testing.T) {
	assert := assert.New(t)

	t.Run("apply", func(t *testing.T) {
		tests := []struct {
			method             string
			uri                string
			expectedStatusCode int
			expectedBody       string
			// mockRaftyErrorMessage error
			file                 string
			createContainerError bool
			raftIsLeader         bool
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
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()

			var payload string
			if tc.file != "" {
				workingDir, err := os.Getwd()
				assert.Nil(err)

				configDir := "testdata/deployments"
				configFile := filepath.Join(workingDir, configDir, tc.file)
				data, err := os.ReadFile(configFile)
				require.NoError(t, err)

				b, err := json.Marshal(APIDeploymentApplyRequest{HCLContent: string(data)})
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
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, header, payload)

			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
		}
	})
}

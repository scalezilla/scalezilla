package cluster

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_api_calls_command(t *testing.T) {
	assert := assert.New(t)

	t.Run("status", func(t *testing.T) {
		tests := []struct {
			makeError  bool
			statusCode int
			response   string
		}{
			{
				statusCode: 200,
				response:   `{"bootstrap": false}`,
			},
			{
				makeError:  true,
				statusCode: 500,
				response:   `{"error", "Internal Server Error"}`,
			},
		}

		for _, tc := range tests {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.makeError {
					_, err := fmt.Fprintln(w, tc.response)
					assert.Nil(err)
					return
				}
				_, err := fmt.Fprintln(w, tc.response)
				assert.Nil(err)
			}))
			defer server.Close()

			config := ClusterHTTPCallBaseConfig{HTTPAddress: server.URL}
			if tc.makeError {
				assert.Error(APICallsBootstrapStatus(config))
			} else {
				assert.Nil(APICallsBootstrapStatus(config))
			}
		}

		// provoke error
		config := ClusterHTTPCallBaseConfig{HTTPAddress: "htttp://127.0.0.1"}
		assert.Error(APICallsBootstrapStatus(config))
	})

	t.Run("bootstrap", func(t *testing.T) {
		tests := []struct {
			makeError  bool
			statusCode int
			token      string
			response   string
		}{
			{
				statusCode: 200,
				response:   `{"accessor_id":"99B82183-8BA0-42E8-B85A-54022781B2C7","token":"F6D7F628-E5B5-4F12-8606-DB6B9E73D5C4","initial_token":true}`,
			},
			{
				statusCode: 200,
				token:      "138762D9-09AC-4A59-A434-79780BC26B19",
				response:   `{"accessor_id":"4C479FE3-9562-41A2-841C-A9A23BCBB844","token":"138762D9-09AC-4A59-A434-79780BC26B19","initial_token":true}`,
			},
			{
				makeError:  true,
				statusCode: 500,
				response:   `{"error", "Internal Server Error"}`,
			},
		}

		for _, tc := range tests {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.makeError {
					_, err := fmt.Fprintln(w, tc.response)
					assert.Nil(err)
					return
				}
				_, err := fmt.Fprintln(w, tc.response)
				assert.Nil(err)
			}))
			defer server.Close()

			config := BootstrapClusterHTTPConfig{Token: tc.token}
			config.HTTPAddress = server.URL
			if tc.makeError {
				assert.Error(APICallsBootstrapCluster(config))
			} else {
				assert.Nil(APICallsBootstrapCluster(config))
			}
		}

		// provoke error
		config := BootstrapClusterHTTPConfig{}
		config.HTTPAddress = "htttp://127.0.0.1"
		assert.Error(APICallsBootstrapCluster(config))
	})

	t.Run("nodes_list", func(t *testing.T) {
		tests := []struct {
			makeError    bool
			statusCode   int
			token        string
			response     string
			kind, format string
		}{
			{
				statusCode: 200,
				response:   `[{"id":"192.168.200.11","name":"server11","address":"192.168.200.11:15002","kind":"server","leader":true,"pool":"default"},{"id":"192.168.200.13","name":"server13","address":"192.168.200.13:15002","kind":"server","leader":false,"pool":"default"},{"id":"192.168.200.12","name":"server12","address":"192.168.200.12:15002","kind":"server","leader":false,"pool":"default"}]`,
				format:     "table",
			},
			{
				statusCode: 200,
				response:   `[{"id":"192.168.200.11","name":"server11","address":"192.168.200.11:15002","kind":"server","leader":true,"pool":"default"},{"id":"192.168.200.13","name":"server13","address":"192.168.200.13:15002","kind":"server","leader":false,"pool":"default"},{"id":"192.168.200.12","name":"server12","address":"192.168.200.12:15002","kind":"server","leader":false,"pool":"default"}]`,
			},
			{
				statusCode: 200,
				kind:       "server",
				response:   `[{"id":"192.168.200.11","name":"server11","address":"192.168.200.11:15002","kind":"server","leader":true,"pool":"default"},{"id":"192.168.200.13","name":"server13","address":"192.168.200.13:15002","kind":"server","leader":false,"pool":"default"},{"id":"192.168.200.12","name":"server12","address":"192.168.200.12:15002","kind":"server","leader":false,"pool":"default"}]`,
			},
			{
				makeError:  true,
				statusCode: 500,
				response:   `{"error", "Internal Server Error"}`,
			},
		}

		for _, tc := range tests {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.makeError {
					_, err := fmt.Fprintln(w, tc.response)
					assert.Nil(err)
					return
				}
				_, err := fmt.Fprintln(w, tc.response)
				assert.Nil(err)
			}))
			defer server.Close()

			format := "json"
			if tc.format != "" {
				format = tc.format
			}
			config := NodesListHTTPConfig{Token: tc.token, Kind: tc.kind}
			config.HTTPAddress = server.URL
			config.OutputFormat = format
			if tc.makeError {
				assert.Error(APICallsNodesList(config))
			} else {
				assert.Nil(APICallsNodesList(config))
			}
		}

		// provoke error
		config := NodesListHTTPConfig{}
		config.HTTPAddress = "htttp://127.0.0.1"
		assert.Error(APICallsNodesList(config))
	})

	t.Run("deployment_apply", func(t *testing.T) {
		tests := []struct {
			makeError, setError bool
			statusCode          int
			token               string
			file                string
			response            string
			format              string
		}{
			{
				statusCode: 200,
				response:   `OK`,
				file:       "testdata/deployments/basic_success.hcl",
			},
			{
				statusCode: 200,
				response:   `OK`,
				file:       "testdata/deployments/basic_success.hcl",
				format:     "json",
			},
			{
				setError:  true,
				makeError: true,
				file:      "testdata/deployments/basic_success.hcl",
				response:  `{"error", "Syntax error"}`,
			},
			{
				makeError: true,
				file:      "testdata/deployments/error_malformed.hcl",
				response:  `{"error", "Syntax error"}`,
			},
			{
				makeError:  true,
				statusCode: 500,
				file:       "testdata/deployments/basic_success.hcl",
				response:   `{"error", "Internal Server Error"}`,
			},
		}

		for _, tc := range tests {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.makeError {
					_, err := fmt.Fprintln(w, tc.response)
					assert.Nil(err)
					return
				}
				_, err := fmt.Fprintln(w, tc.response)
				assert.Nil(err)
			}))
			defer server.Close()

			config := DeploymentApplyHTTPConfig{
				Token: tc.token,
				File:  tc.file,
			}
			config.HTTPAddress = server.URL
			config.OutputFormat = tc.format
			if tc.file != "" {
				workingDir, err := os.Getwd()
				assert.Nil(err)
				config.File = filepath.Join(workingDir, tc.file)
			}

			if tc.setError {
				config.osReadFile = func(name string) ([]byte, error) {
					return nil, fmt.Errorf("os read file error")
				}
			}

			if tc.makeError {
				assert.Error(APICallsDeploymentApply(config))
			} else {
				assert.Nil(APICallsDeploymentApply(config))
			}
		}

		// provoke error
		config := DeploymentApplyHTTPConfig{}
		workingDir, err := os.Getwd()
		assert.Nil(err)
		config.File = filepath.Join(workingDir, "testdata/deployments/basic_success.hcl")
		config.HTTPAddress = "htttp://127.0.0.1"
		assert.Error(APICallsDeploymentApply(config))
	})

	t.Run("pods_list", func(t *testing.T) {
		tests := []struct {
			makeError, setError bool
			statusCode          int
			token               string
			response            string
			format              string
		}{
			{
				statusCode: 200,
				response:   `[{"namespace":"default","id":"nginx-test-nginx-container","image":"docker.io/library/nginx:latest","pid":5086,"runtime":"io.containerd.runc.v2","status":"RUNNING","created_at":"0001-01-01T00:00:00Z"}]`,
			},
			{
				statusCode: 200,
				response:   `[{"namespace":"default","id":"nginx-test-nginx-container","image":"docker.io/library/nginx:latest","pid":5086,"runtime":"io.containerd.runc.v2","status":"RUNNING","created_at":"0001-01-01T00:00:00Z"}]`,
				format:     "json",
			},
			{
				makeError:  true,
				statusCode: 500,
				response:   `{"error", "Internal Server Error"}`,
			},
		}

		for _, tc := range tests {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				if tc.makeError {
					_, err := fmt.Fprintln(w, tc.response)
					assert.Nil(err)
					return
				}
				_, err := fmt.Fprintln(w, tc.response)
				assert.Nil(err)
			}))
			defer server.Close()

			namespace := "default"
			config := PodsListHTTPConfig{
				Token:     tc.token,
				Namespace: namespace,
			}
			config.HTTPAddress = server.URL
			config.OutputFormat = tc.format

			if tc.makeError {
				assert.Error(APICallsPodsList(config))
			} else {
				assert.Nil(APICallsPodsList(config))
			}
		}

		// provoke error
		config := PodsListHTTPConfig{}
		config.HTTPAddress = "htttp://127.0.0.1"
		assert.Error(APICallsPodsList(config))
	})
}

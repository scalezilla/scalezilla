package cluster

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
			APICallsBootstrapStatus(config)
		}

		// provoke error
		config := ClusterHTTPCallBaseConfig{HTTPAddress: "htttp://127.0.0.1"}
		APICallsBootstrapStatus(config)
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
			APICallsBootstrapCluster(config)
		}

		// provoke error
		config := BootstrapClusterHTTPConfig{}
		config.HTTPAddress = "htttp://127.0.0.1"
		APICallsBootstrapCluster(config)
	})
}

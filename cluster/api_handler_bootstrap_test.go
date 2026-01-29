package cluster

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_api_handlers_bootstrap(t *testing.T) {
	assert := assert.New(t)

	t.Run("bootstrap", func(t *testing.T) {
		tests := []struct {
			method                string
			uri                   string
			expectedStatusCode    int
			expectedBody          string
			mockRaftyErrorMessage error
			bootstrapped          bool
			header                map[string]string
			body                  string
			errorSubmitCommand    bool
			errorBootstrapCluster bool
		}{
			{
				method:             "GET",
				uri:                "/api/v1/cluster/bootstrap/status",
				expectedStatusCode: 200,
				expectedBody:       `{"bootstrapped":false}`,
			},
			{
				method:             "POST",
				uri:                "/api/v1/cluster/bootstrap/cluster",
				expectedStatusCode: 403,
				expectedBody:       `{"error":"cluster already boostrapped"}`,
				bootstrapped:       true,
				header: map[string]string{
					"Content-Type": "application/json",
				},
			},
			{
				method:             "POST",
				uri:                "/api/v1/cluster/bootstrap/cluster",
				expectedStatusCode: 200,
				expectedBody:       `{"accessor_id":"`,
				header: map[string]string{
					"Content-Type": "application/json",
				},
			},
			{
				method:             "POST",
				uri:                "/api/v1/cluster/bootstrap/cluster",
				expectedStatusCode: 400,
				expectedBody:       `{"error":"`,
				header: map[string]string{
					"Content-Type": "application/json",
				},
				body: `{`,
			},
			{
				method:             "POST",
				uri:                "/api/v1/cluster/bootstrap/cluster",
				expectedStatusCode: 200,
				expectedBody:       `{"accessor_id":"`,
				header: map[string]string{
					"Content-Type": "application/json",
				},
				body: `{}`,
			},
			{
				method:             "POST",
				uri:                "/api/v1/cluster/bootstrap/cluster",
				expectedStatusCode: 200,
				expectedBody:       `{"accessor_id":"`,
				header: map[string]string{
					"Content-Type": "application/json",
				},
				body: `{"token":"OK"}`,
			},
			{
				method:             "POST",
				uri:                "/api/v1/cluster/bootstrap/cluster",
				expectedStatusCode: 500,
				expectedBody:       `{"error":"`,
				header: map[string]string{
					"Content-Type": "application/json",
				},
				body:                  `{}`,
				errorBootstrapCluster: true,
			},
			{
				method:             "POST",
				uri:                "/api/v1/cluster/bootstrap/cluster",
				expectedStatusCode: 500,
				expectedBody:       `{"error":"`,
				header: map[string]string{
					"Content-Type": "application/json",
				},
				body:               `{}`,
				errorSubmitCommand: true,
			},
		}

		for _, tc := range tests {
			cfg := basicClusterConfig{randomPort: false, dev: true}
			cluster := makeBasicCluster(cfg)
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()
			mock := mockRafty{}
			cluster.rafty = &mock
			mock.bootstrapped = tc.bootstrapped
			if tc.errorBootstrapCluster {
				mock.errBootstrap = errors.New("bootstrap error")
			}
			if tc.errorSubmitCommand {
				mock.err = errors.New("submit error")
			}
			router := cluster.newApiRouters()
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, tc.header, tc.body)

			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
			assert.Contains(w.Body.String(), tc.expectedBody, "Failed to get right body content")
			mock.err = nil
		}
	})
}

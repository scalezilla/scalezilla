package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIHandlers(t *testing.T) {
	assert := assert.New(t)

	t.Run("health", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		router := cluster.newApiRouters()

		tests := []struct {
			method             string
			uri                string
			expectedStatusCode int
			expectedBody       string
		}{
			{
				method:             "GET",
				uri:                "/api/v1/cluster/health",
				expectedStatusCode: 200,
				expectedBody:       `{"message":"OK"}`,
			},
			{
				method:             "GET",
				uri:                "/api/v1/cluster/healthz",
				expectedStatusCode: 200,
				expectedBody:       `{"message":"OK"}`,
			},
		}

		for _, tc := range tests {
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, nil)

			assert.Equal(200, w.Code, "Failed to perform http GET request")
			assert.Contains(w.Body.String(), tc.expectedBody, "Failed to get right body content")
		}
	})
}

package cluster

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/scalezilla/scalezilla/cri"
	"github.com/stretchr/testify/assert"
)

func TestCluster_api_handler_pods(t *testing.T) {
	assert := assert.New(t)

	t.Run("list", func(t *testing.T) {
		tests := []struct {
			method             string
			uri                string
			expectedStatusCode int
			expectedBody       string
			listContainerError bool
			bootstrapped       bool
			raftIsLeader       bool
			raftyLeader        bool
			result             []cri.ContainerList
		}{
			{
				method:             "GET",
				uri:                "/api/v1/pods/list?namespace=default",
				expectedStatusCode: 200,
				expectedBody:       `[{"namespace":"default","id":"nginx-test-nginx-container","image":"docker.io/library/nginx:latest","pid":5086,"runtime":"io.containerd.runc.v2","status":"RUNNING","created_at":"0001-01-01T00:00:00Z"}]`,
				bootstrapped:       true,
				raftyLeader:        true,
				result: []cri.ContainerList{
					{
						Namespace: "default",
						ID:        "nginx-test-nginx-container",
						PID:       5086,
						Image:     "docker.io/library/nginx:latest",
						Runtime:   "io.containerd.runc.v2",
						Status:    "RUNNING",
					},
				},
			},
			{
				method:             "GET",
				uri:                "/api/v1/pods/list?namespace=all",
				expectedStatusCode: 200,
				expectedBody:       `[{"namespace":"default","id":"nginx-test-nginx-container","image":"docker.io/library/nginx:latest","pid":5086,"runtime":"io.containerd.runc.v2","status":"RUNNING","created_at":"0001-01-01T00:00:00Z"}]`,
				bootstrapped:       true,
				raftyLeader:        true,
				result: []cri.ContainerList{
					{
						Namespace: "default",
						ID:        "nginx-test-nginx-container",
						PID:       5086,
						Image:     "docker.io/library/nginx:latest",
						Runtime:   "io.containerd.runc.v2",
						Status:    "RUNNING",
					},
				},
			},
			{
				method:             "GET",
				uri:                "/api/v1/pods/list?namespace=plop",
				expectedStatusCode: 404,
				expectedBody:       `[]`,
				bootstrapped:       true,
				raftyLeader:        true,
			},
			{
				method:             "GET",
				uri:                "/api/v1/pods/list",
				expectedStatusCode: 400,
				bootstrapped:       true,
				raftyLeader:        true,
				listContainerError: true,
			},
			{
				method:             "GET",
				uri:                "/api/v1/pods/list",
				expectedStatusCode: 403,
				expectedBody:       `{"error":"cluster not boostrapped"}`,
			},
			{
				method:             "GET",
				uri:                "/api/v1/pods/list",
				expectedStatusCode: 403,
				expectedBody:       `{"error":"no leader"}`,
				bootstrapped:       true,
			},
		}

		for _, tc := range tests {
			cfg := basicClusterConfig{randomPort: false, dev: true}
			cluster := makeBasicCluster(cfg)
			defer func() {
				_ = os.RemoveAll(cluster.config.DataDir)
			}()

			mock := mockRafty{}
			if tc.raftyLeader {
				mock.isLeader = true
			}
			cluster.rafty = &mock

			if tc.listContainerError {
				cluster.di.listContainerFunc = func(ctx context.Context, namespace string) ([]cri.ContainerList, error) {
					return nil, errors.New("container error")
				}
			} else {
				cluster.di.listContainerFunc = func(ctx context.Context, namespace string) ([]cri.ContainerList, error) {
					if tc.expectedStatusCode == 200 {
						return tc.result, nil
					}
					return nil, nil
				}
			}

			router := cluster.newApiRouters()
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, nil, "")

			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
			switch tc.expectedStatusCode {
			case 404:
			case 200:
				assert.Equal(tc.expectedBody, w.Body.String(), fmt.Sprintf("Failed to perform http %s request", tc.method))

			}
		}
	})
}

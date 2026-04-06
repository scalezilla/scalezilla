package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/scalezilla/scalezilla/cri"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("delete", func(t *testing.T) {
		tests := []struct {
			method               string
			uri                  string
			expectedStatusCode   int
			expectedBody         string
			deleteContainerError bool
			bootstrapped         bool
			raftIsLeader         bool
			raftyLeader          bool
			pods                 []string
			namespace            string
			detached             bool
		}{
			{
				method:             "DELETE",
				uri:                "/api/v1/pods/delete",
				expectedStatusCode: 200,
				pods:               []string{"nginx-test-nginx-container"},
				bootstrapped:       true,
				raftyLeader:        true,
			},
			{
				method:             "DELETE",
				uri:                "/api/v1/pods/delete",
				expectedStatusCode: 200,
				pods:               []string{"nginx-test-nginx-container"},
				bootstrapped:       true,
				raftyLeader:        true,
				namespace:          "default",
			},
			{
				method:             "DELETE",
				uri:                "/api/v1/pods/delete",
				expectedStatusCode: 200,
				pods:               []string{"nginx-test-nginx-container"},
				bootstrapped:       true,
				raftyLeader:        true,
				namespace:          "default",
				detached:           true,
			},
			{
				method:               "DELETE",
				uri:                  "/api/v1/pods/delete",
				pods:                 []string{"nginx-test-nginx-container"},
				expectedStatusCode:   200,
				bootstrapped:         true,
				raftyLeader:          true,
				deleteContainerError: true,
			},
			{
				method:               "DELETE",
				uri:                  "/api/v1/pods/delete",
				pods:                 []string{"nginx-test-nginx-container"},
				expectedStatusCode:   200,
				bootstrapped:         true,
				raftyLeader:          true,
				deleteContainerError: true,
				detached:             true,
			},
			{
				method:             "DELETE",
				uri:                "/api/v1/pods/delete",
				expectedStatusCode: 400,
				bootstrapped:       true,
				raftyLeader:        true,
			},
			{
				method:             "DELETE",
				uri:                "/api/v1/pods/delete",
				pods:               []string{"nginx-test-nginx-container"},
				expectedStatusCode: 403,
				expectedBody:       `{"error":"cluster not boostrapped"}`,
			},
			{
				method:             "DELETE",
				uri:                "/api/v1/pods/delete",
				pods:               []string{"nginx-test-nginx-container"},
				expectedStatusCode: 403,
				expectedBody:       `{"error":"no leader"}`,
				bootstrapped:       true,
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
			if len(tc.pods) > 0 {
				b, err := json.Marshal(APIPodsDeleteRequest{
					Namespace: tc.namespace,
					Pods:      tc.pods,
					Detached:  tc.detached,
				})
				require.NoError(t, err)
				payload = string(b)
			}

			mock := mockRafty{}
			if tc.raftyLeader {
				mock.isLeader = true
			}
			cluster.rafty = &mock

			if tc.deleteContainerError {
				cluster.di.deleteContainerFunc = func(ctx context.Context, namespace, containerID string, stopTimeout time.Duration) error {
					return errors.New("container error")
				}
			} else {
				cluster.di.deleteContainerFunc = func(ctx context.Context, namespace, containerID string, stopTimeout time.Duration) error {
					return nil
				}
			}

			router := cluster.newApiRouters()
			w := makeHTTPRequestRecorder(router, tc.method, tc.uri, header, payload)

			assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
			// switch tc.expectedStatusCode {
			// case 404:
			// case 200:
			// 	assert.Equal(tc.expectedBody, w.Body.String(), fmt.Sprintf("Failed to perform http %s request", tc.method))
			// }
		}
	})
}

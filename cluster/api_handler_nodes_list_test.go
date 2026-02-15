package cluster

import (
	"fmt"
	"os"
	"testing"

	"github.com/Lord-Y/rafty"
	"github.com/scalezilla/scalezilla/osdiscovery"
	"github.com/stretchr/testify/assert"
)

func TestCluster_api_handler_nodes_list(t *testing.T) {
	assert := assert.New(t)

	t.Run("list", func(t *testing.T) {
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
			dev                   bool
			raftyStatus           rafty.Status
			raftyLeader           bool
			clientSize            uint16
		}{
			{
				method:             "GET",
				uri:                "/api/v1/cluster/nodes/list",
				expectedStatusCode: 200,
				expectedBody:       `[{"id":"127.0.0.1:15002","name":"dev","address":"127.0.0.1:15002","kind":"server","leader":true,"nodePool":"default"}]`,
				dev:                true,
			},
			{
				method:             "GET",
				uri:                "/api/v1/cluster/nodes/list",
				expectedStatusCode: 403,
				expectedBody:       `{"error":"cluster not boostrapped"}`,
			},
			{
				method:             "GET",
				uri:                "/api/v1/cluster/nodes/list",
				expectedStatusCode: 403,
				expectedBody:       `{"error":"no leader"}`,
				bootstrapped:       true,
			},
			{
				method:             "GET",
				uri:                "/api/v1/cluster/nodes/list",
				expectedStatusCode: 200,
				bootstrapped:       true,
				raftyLeader:        true,
				expectedBody:       `[{"id":"127.0.0.1:20008","name":"fake","address":"127.0.0.1:20008","kind":"server","leader":true,"nodePool":"default"},{"id":"127.0.0.1:20001","name":"127.0.0.1:20001","address":"127.0.0.1:20001","kind":"server","leader":false,"nodePool":"default"},{"id":"127.0.0.1:20004","name":"127.0.0.1:20004","address":"127.0.0.1:20004","kind":"server","leader":false,"nodePool":"default"}]`,
			},
			{
				method:             "GET",
				uri:                "/api/v1/cluster/nodes/list?kind=server",
				expectedStatusCode: 200,
				bootstrapped:       true,
				raftyLeader:        true,
				expectedBody:       `[{"id":"127.0.0.1:20008","name":"fake","address":"127.0.0.1:20008","kind":"server","leader":true,"nodePool":"default"},{"id":"127.0.0.1:20001","name":"127.0.0.1:20001","address":"127.0.0.1:20001","kind":"server","leader":false,"nodePool":"default"},{"id":"127.0.0.1:20004","name":"127.0.0.1:20004","address":"127.0.0.1:20004","kind":"server","leader":false,"nodePool":"default"}]`,
			},
			{
				method:             "GET",
				uri:                "/api/v1/cluster/nodes/list",
				expectedStatusCode: 200,
				bootstrapped:       true,
				raftyLeader:        true,
				clientSize:         1,
				expectedBody:       `[{"id":"127.0.0.1:20001","name":"127.0.0.1:20001","address":"127.0.0.1:20001","kind":"server","leader":true,"nodePool":"default"},{"id":"127.0.0.1:20004","name":"127.0.0.1:20004","address":"127.0.0.1:20004","kind":"server","leader":false,"nodePool":"default"},{"id":"127.0.0.1:20007","name":"127.0.0.1:20007","address":"127.0.0.1:20007","kind":"server","leader":false,"nodePool":"default"},{"id":"127.0.0.1:20005","name":"fake","address":"127.0.0.1:20005","kind":"client","leader":false,"nodePool":"default"}]`,
			},
			{
				method:             "GET",
				uri:                "/api/v1/cluster/nodes/list?kind=client",
				expectedStatusCode: 200,
				bootstrapped:       true,
				raftyLeader:        true,
				clientSize:         1,
				expectedBody:       `[{"id":"127.0.0.1:20005","name":"fake","address":"127.0.0.1:20005","kind":"client","leader":false,"nodePool":"default"}]`,
			},
		}

		for _, tc := range tests {
			if tc.dev {
				cfg := basicClusterConfig{randomPort: false, dev: tc.dev}
				cluster := makeBasicCluster(cfg)
				defer func() {
					_ = os.RemoveAll(cluster.config.DataDir)
				}()
				assert.Nil(cluster.checkSystemInfo())
				// the following is only a fix only for CI
				cluster.systemInfo.OS.Hostname = "dev"
				router := cluster.newApiRouters()
				w := makeHTTPRequestRecorder(router, tc.method, tc.uri, tc.header, tc.body)

				assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
				assert.Contains(w.Body.String(), tc.expectedBody, "Failed to get right body content")
			} else {
				clusters := makeSizedCluster(sizedClusterConfig{clientSize: tc.clientSize})
				cluster := clusters[len(clusters)-1]
				defer func() {
					_ = os.RemoveAll(cluster.config.DataDir)
				}()

				mock := mockRafty{}
				cluster.rafty = &mock
				mock.bootstrapped = tc.bootstrapped
				mock.leader = tc.raftyLeader

				tc.raftyStatus = rafty.Status{Configuration: rafty.Configuration{}}
				if tc.clientSize > 0 {
					if tc.raftyLeader {
						mock.leaderAddress = cluster.config.Client.ClusterJoin.InitialMembers[0]
						mock.leaderId = cluster.config.Client.ClusterJoin.InitialMembers[0]
					}
					for i, member := range cluster.config.Client.ClusterJoin.InitialMembers {
						var voter bool
						if i >= int(tc.clientSize) {
							voter = true
						}
						cluster.nodeMap[member] = &nodeMap{
							ID:       member,
							Address:  member,
							IsVoter:  voter,
							NodePool: defaultNodePool,
							SystemInfo: osdiscovery.SystemInfo{
								OS: &osdiscovery.OS{
									Hostname: member,
								},
							},
						}
					}

					for i, x := range cluster.buildPeers() {
						var voter bool
						if i >= int(tc.clientSize) {
							voter = true
						}
						tc.raftyStatus.Configuration.ServerMembers = append(tc.raftyStatus.Configuration.ServerMembers, rafty.Peer{
							Address: x.Address,
							ID:      x.Address,
							IsVoter: voter,
						})
					}
				} else {
					if tc.raftyLeader {
						mock.leaderAddress = cluster.raftyAddress.String()
						mock.leaderId = cluster.raftyAddress.String()
					}

					for _, member := range cluster.config.Server.ClusterJoin.InitialMembers {
						cluster.nodeMap[member] = &nodeMap{
							ID:       member,
							Address:  member,
							IsVoter:  true,
							NodePool: defaultNodePool,
							SystemInfo: osdiscovery.SystemInfo{
								OS: &osdiscovery.OS{
									Hostname: member,
								},
							},
						}
					}

					for _, x := range cluster.buildPeers() {
						tc.raftyStatus.Configuration.ServerMembers = append(tc.raftyStatus.Configuration.ServerMembers, rafty.Peer{
							Address: x.Address,
							ID:      x.Address,
							IsVoter: true,
						})
					}
				}

				mock.raftyStatus = tc.raftyStatus
				assert.Nil(cluster.checkSystemInfo())
				cluster.config.Hostname = "fake"

				router := cluster.newApiRouters()
				w := makeHTTPRequestRecorder(router, tc.method, tc.uri, tc.header, tc.body)

				assert.Equal(tc.expectedStatusCode, w.Code, fmt.Sprintf("Failed to perform http %s request", tc.method))
				assert.Contains(w.Body.String(), tc.expectedBody, "Failed to get right body content")
				mock.err = nil
			}
		}
	})
}

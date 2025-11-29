package cluster

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_responses(t *testing.T) {
	assert := assert.New(t)

	t.Run("resp_service_ports_discovery", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()

		resp := RPCResponse{
			Error: errors.New("test error"),
		}
		cluster.respServicePortsDiscovery(resp)

		resp.Error = nil
		response := RPCServicePortsDiscoveryResponse{}
		resp.Response = response
		cluster.respServicePortsDiscovery(resp)

		response = RPCServicePortsDiscoveryResponse{
			IsVoter: true,
			Address: "1234",
			ID:      "test",
		}
		resp.Response = response
		cluster.respServicePortsDiscovery(resp)
		assert.Equal(cluster.nodeMap[response.ID].Address, response.Address)
		assert.Equal(cluster.nodeMap[response.ID].IsVoter, response.IsVoter)
	})
}

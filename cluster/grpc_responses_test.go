package cluster

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_responses(t *testing.T) {
	assert := assert.New(t)

	t.Run("resp_service_ports_discovery", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.buildAddressAndID()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

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

	t.Run("resp_service_node_polling", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.buildAddressAndID()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(cluster.checkSystemInfo())

		resp := RPCResponse{
			Error: errors.New("test error"),
		}
		cluster.respServiceNodePolling(resp)

		resp.Error = nil
		response := RPCServiceNodePollingResponse{}
		resp.Response = response
		cluster.respServiceNodePolling(resp)

		response = RPCServiceNodePollingResponse{
			Address: "1234",
			ID:      "test",
		}
		resp.Response = response
		cluster.respServiceNodePolling(resp)
		assert.Equal(cluster.nodeMap[response.ID].Address, response.Address)
		assert.Equal(cluster.nodeMap[response.ID].ID, response.ID)
	})
}

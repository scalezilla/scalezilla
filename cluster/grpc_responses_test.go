package cluster

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_responses(t *testing.T) {
	assert := assert.New(t)

	t.Run("resp_service_ports_discovery_server", func(t *testing.T) {
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

	t.Run("resp_service_ports_discovery_client", func(t *testing.T) {
		clusters := makeSizedCluster(sizedClusterConfig{clientSize: 1})
		cluster := clusters[len(clusters)-1]
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

		cluster.nodeMap["1234"] = &nodeMap{
			Address: "1234",
			ID:      "1234",
		}

		response = RPCServiceNodePollingResponse{
			Address:        "1234",
			ID:             "1234",
			OsHostname:     cluster.systemInfo.OS.Hostname,
			OsArchitecture: cluster.systemInfo.OS.Architecture,
		}
		resp.Response = response
		cluster.respServiceNodePolling(resp)
		assert.NotEmpty(cluster.nodeMap[response.ID].SystemInfo.OS.Hostname)
		assert.NotEmpty(cluster.nodeMap[response.ID].SystemInfo.OS.Architecture)
	})
}

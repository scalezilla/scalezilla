package cluster

import (
	"os"
	"testing"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_receiver_responses(t *testing.T) {
	assert := assert.New(t)

	t.Run("rcv_service_ports_discovery", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.buildAddressAndID()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		responseChan := make(chan RPCResponse, 1)
		request := RPCRequest{
			RPCType: ServicePortsDiscovery,
			Request: &scalezillapb.ServicePortsDiscoveryRequestReply{
				Address:  "1234",
				Id:       "1234",
				PortHttp: uint32(defaultHTTPPort),
				PortGrpc: uint32(defaultGRPCPort),
				PortRaft: uint32(defaultRaftGRPCPort),
				IsVoter:  cluster.isVoter,
				NodePool: cluster.nodePool,
			},
			ResponseChan: responseChan,
		}

		go cluster.rcvServicePortsDiscovery(request)
		response := <-responseChan
		assert.NotNil(response.Response)
	})

	t.Run("rcv_service_node_polling", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.buildAddressAndID()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		assert.Nil(cluster.checkSystemInfo())
		responseChan := make(chan RPCResponse, 1)
		request := RPCRequest{
			RPCType: ServiceNodePolling,
			Request: &scalezillapb.ServiceNodePollingRequestReply{
				Address:                "1234",
				Id:                     "1234",
				OsName:                 cluster.systemInfo.OS.Name,
				OsVendor:               cluster.systemInfo.OS.Vendor,
				OsVersion:              cluster.systemInfo.OS.Version,
				OsFamily:               cluster.systemInfo.OS.Family,
				OsHostname:             cluster.systemInfo.OS.Hostname,
				OsArchitecture:         cluster.systemInfo.OS.Architecture,
				OsType:                 cluster.systemInfo.OS.OSType,
				CpuTotal:               cluster.systemInfo.CPU.CPU,
				CpuCores:               cluster.systemInfo.CPU.Cores,
				CpuFrequency:           cluster.systemInfo.CPU.Frequency,
				CpuCumulativeFrequency: cluster.systemInfo.CPU.CumulativeFrequency,
				CpuCapabilitites:       cluster.systemInfo.CPU.Capabilitites,
				CpuVendor:              cluster.systemInfo.CPU.Vendor,
				CpuModel:               cluster.systemInfo.CPU.Model,
				MemoryTotal:            uint64(cluster.systemInfo.Memory.Total),
				MemoryAvailable:        uint64(cluster.systemInfo.Memory.Available),
			},
			ResponseChan: responseChan,
		}

		go cluster.rcvServiceNodePolling(request)
		response := <-responseChan
		assert.NotNil(response.Response)
	})
}

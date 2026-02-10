package cluster

import (
	"os"
	"testing"
	"time"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_requests(t *testing.T) {
	assert := assert.New(t)

	t.Run("send_rpc_service_ports_discovery", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		request := RPCRequest{
			RPCType: ServicePortsDiscovery,
			Request: RPCServicePortsDiscoveryRequest{
				Address:  cluster.grpcAddress.String(),
				ID:       cluster.id,
				NodePool: cluster.nodePool,
				PortHTTP: uint32(cluster.config.HTTPPort),
				PortGRPC: uint32(cluster.config.GRPCPort),
				PortRaft: uint32(cluster.config.RaftGRPCPort),
				IsVoter:  cluster.isVoter,
			},
			Timeout:      time.Second,
			ResponseChan: cluster.rpcServicePortsDiscoveryChanResp,
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			data := <-cluster.rpcServicePortsDiscoveryChanResp
			assert.ErrorContains(data.Error, "connection refused")
		}()

		member := cluster.members[0]
		if client := cluster.getClient(member); client != nil {
			cluster.di.sendRPCFunc(member, client, request)
		}
	})

	t.Run("req_service_ports_discovery", func(t *testing.T) {
		clusters := makeSizedCluster(sizedClusterConfig{})
		cluster := clusters[0]

		resp := &scalezillapb.ServicePortsDiscoveryRequestReply{
			Address:  "1234",
			Id:       "1234",
			PortHttp: uint32(defaultHTTPPort),
			PortGrpc: uint32(defaultGRPCPort),
			PortRaft: uint32(defaultRaftGRPCPort),
			IsVoter:  cluster.isVoter,
			NodePool: cluster.nodePool,
		}
		cluster.di.sendRPCFunc = func(address string, client scalezillapb.ScalezillaClient, request RPCRequest) {
			request.ResponseChan <- RPCResponse{
				Response: makeServicePortsDiscoveryResponse(resp), Error: nil, TargetNode: address}
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			data := <-cluster.rpcServicePortsDiscoveryChanResp
			response := data.Response.(RPCServicePortsDiscoveryResponse)
			assert.Equal(resp.Address, response.Address)
		}()

		cluster.reqServicePortsDiscovery()
	})

	t.Run("send_rpc_service_node_polling", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		assert.Nil(cluster.checkSystemInfo())
		request := RPCRequest{
			RPCType: ServiceNodePolling,
			Request: RPCServiceNodePollingRequest{
				Address:                cluster.grpcAddress.String(),
				ID:                     cluster.id,
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
			Timeout:      time.Second,
			ResponseChan: cluster.rpcServiceNodePollingChanResp,
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			data := <-cluster.rpcServiceNodePollingChanResp
			assert.ErrorContains(data.Error, "connection refused")
		}()

		member := cluster.members[0]
		if client := cluster.getClient(member); client != nil {
			cluster.di.sendRPCFunc(member, client, request)
		}
	})

	t.Run("req_service_node_polling", func(t *testing.T) {
		clusters := makeSizedCluster(sizedClusterConfig{})
		cluster := clusters[0]
		assert.Nil(cluster.checkSystemInfo())

		resp := &scalezillapb.ServiceNodePollingRequestReply{
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
		}
		cluster.di.sendRPCFunc = func(address string, client scalezillapb.ScalezillaClient, request RPCRequest) {
			request.ResponseChan <- RPCResponse{
				Response: makeServiceNodePollingResponse(resp), Error: nil, TargetNode: address}
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			data := <-cluster.rpcServiceNodePollingChanResp
			response := data.Response.(RPCServiceNodePollingResponse)
			assert.Equal(resp.Address, response.Address)
		}()

		cluster.reqServiceNodePolling()
	})
}

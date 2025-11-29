package cluster

import (
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
}

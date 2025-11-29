package cluster

import (
	"context"
	"testing"
	"time"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_loop(t *testing.T) {
	assert := assert.New(t)

	t.Run("grpc_loop", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		ctx, cancel := context.WithCancel(context.Background())
		cluster.ctx = ctx
		cluster.wg.Go(cluster.grpcLoop)

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

		go func() {
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		cluster.rpcServicePortsDiscoveryChanReq <- request
		response := <-responseChan
		assert.NotNil(response.Response)

		cluster.rpcServicePortsDiscoveryChanResp <- RPCResponse{
			Response: RPCServicePortsDiscoveryResponse{
				IsVoter: true,
				Address: "1234",
				ID:      "test",
			},
		}

		cluster.wg.Wait()
	})

	t.Run("drain_rcv_service_ports_discovery", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)

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

		cluster.wg.Go(func() {
			for {
				select {
				case cluster.rpcServicePortsDiscoveryChanReq <- request:
				case <-time.After(200 * time.Millisecond):
					return
				}
			}
		})

		cluster.wg.Go(func() {
			time.Sleep(100 * time.Millisecond)
			cluster.drainRCVServicePortsDiscovery()
		})
		cluster.wg.Wait()
	})

	t.Run("drain_resp_service_ports_discovery", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)

		cluster.wg.Go(func() {
			for {
				select {
				case cluster.rpcServicePortsDiscoveryChanResp <- RPCResponse{}:
				case <-time.After(200 * time.Millisecond):
					return
				}
			}
		})

		cluster.wg.Go(func() {
			time.Sleep(100 * time.Millisecond)
			cluster.drainRespServicePortsDiscovery()
		})
		cluster.wg.Wait()
	})
}

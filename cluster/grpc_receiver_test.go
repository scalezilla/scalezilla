package cluster

import (
	"context"
	"testing"
	"time"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_receiver(t *testing.T) {
	assert := assert.New(t)

	t.Run("service_ports_discovery_context_cancelled_first", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.ctx = context.Background()
		cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()

		request := &scalezillapb.ServicePortsDiscoveryRequestReply{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		_, err := cluster.ServicePortsDiscovery(ctx, request)
		assert.ErrorIs(context.Canceled, err)
	})

	t.Run("service_ports_discovery_context_shutdown_first", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		ctx, cancel := context.WithCancel(context.Background())
		cluster.ctx = ctx
		cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()
		request := &scalezillapb.ServicePortsDiscoveryRequestReply{}

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		_, err := cluster.ServicePortsDiscovery(context.Background(), request)
		assert.ErrorIs(ErrShutdown, err)
	})

	t.Run("service_ports_discovery_err_timeout_first", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.ctx = context.Background()
		cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()

		request := &scalezillapb.ServicePortsDiscoveryRequestReply{}

		go func() {
			time.Sleep(100 * time.Millisecond)
			cluster.ctx.Done()
		}()

		_, err := cluster.ServicePortsDiscovery(context.Background(), request)
		assert.ErrorIs(ErrTimeout, err)
	})

	t.Run("service_ports_discovery_response", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.ctx = context.Background()
		cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()

		request := &scalezillapb.ServicePortsDiscoveryRequestReply{
			Address:  "12345",
			Id:       "12345",
			PortHttp: uint32(defaultHTTPPort),
			PortGrpc: uint32(defaultGRPCPort),
			PortRaft: uint32(defaultRaftGRPCPort),
			IsVoter:  true,
			NodePool: defaultNodePool,
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			data := <-cluster.rpcServicePortsDiscoveryChanReq
			data.ResponseChan <- RPCResponse{
				Response: &scalezillapb.ServicePortsDiscoveryRequestReply{
					Address:  "123456",
					Id:       "123456",
					PortHttp: uint32(defaultHTTPPort),
					PortGrpc: uint32(defaultGRPCPort),
					PortRaft: uint32(defaultRaftGRPCPort),
					IsVoter:  true,
					NodePool: defaultNodePool,
				},
			}
		}()

		response, err := cluster.ServicePortsDiscovery(context.Background(), request)
		assert.Nil(err)
		assert.Equal("123456", response.Id)
	})

	t.Run("service_ports_discovery_context_cancelled_second", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.ctx = context.Background()
		cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()

		request := &scalezillapb.ServicePortsDiscoveryRequestReply{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			<-cluster.rpcServicePortsDiscoveryChanReq
		}()

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		_, err := cluster.ServicePortsDiscovery(ctx, request)
		assert.ErrorIs(context.Canceled, err)
	})

	t.Run("service_ports_discovery_context_shutdown_second", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		ctx, cancel := context.WithCancel(context.Background())
		cluster.ctx = ctx
		// cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()
		request := &scalezillapb.ServicePortsDiscoveryRequestReply{}

		go func() {
			<-cluster.rpcServicePortsDiscoveryChanReq
		}()

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		_, err := cluster.ServicePortsDiscovery(context.Background(), request)
		assert.ErrorIs(ErrShutdown, err)
	})

	t.Run("service_ports_discovery_err_timeout_second", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.ctx = context.Background()
		cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()

		request := &scalezillapb.ServicePortsDiscoveryRequestReply{}

		go func() {
			<-cluster.rpcServicePortsDiscoveryChanReq
			time.Sleep(time.Second)
		}()

		_, err := cluster.ServicePortsDiscovery(context.Background(), request)
		assert.ErrorIs(ErrTimeout, err)
	})
}

package cluster

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_server(t *testing.T) {
	assert := assert.New(t)

	t.Run("start_grpc_server_listen_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		cluster.di.grpcListenFunc = func(network, address string) (net.Listener, error) {
			return nil, errors.New("start error")
		}

		assert.Error(cluster.startGRPCServer())
	})

	t.Run("start_grpc_server_serve_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		cluster.di.grpcListenFunc = func(network, address string) (net.Listener, error) {
			return badListener{}, nil
		}
		cluster.grpcAddress = net.TCPAddr{}

		assert.Error(cluster.startGRPCServer())
	})

	t.Run("start_stop_grpc_server_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.buildAddressAndID()
		ctx, cancel := context.WithCancel(context.Background())
		cluster.ctx = ctx
		defer func() {
			cancel()
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		go func() {
			assert.Nil(cluster.startGRPCServer())
		}()

		time.Sleep(100 * time.Millisecond)
		cluster.stopGRPCServer()
	})

	t.Run("get_client", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.buildAddressAndID()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		assert.NotNil(cluster.getClient(cluster.members[0]))
		assert.NotNil(cluster.getClient(cluster.members[0])) // second time to fetch data from map
	})

	t.Run("cluster_dev", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.dev = true
		cluster.checkBootstrapSize()
	})

	t.Run("check_bootstrap_size_context_done", func(t *testing.T) {
		clusters := makeSizedCluster(sizedClusterConfig{})
		cluster := clusters[0]
		ctx, cancel := context.WithCancel(context.Background())
		cluster.ctx = ctx
		cancel()
		go cluster.checkBootstrapSize()
		cluster.wg.Wait()
	})

	t.Run("check_bootstrap_size_reached", func(t *testing.T) {
		clusters := makeSizedCluster(sizedClusterConfig{})
		cluster := clusters[0]
		ctx, cancel := context.WithCancel(context.Background())
		cluster.ctx = ctx
		cluster.checkBootstrapSizeDuration = 500 * time.Millisecond

		cluster.wg.Go(func() {
			for {
				select {
				case <-cluster.ctx.Done():
					return

				case data, ok := <-cluster.rpcServicePortsDiscoveryChanReq:
					if ok {
						cluster.rcvServicePortsDiscovery(data)
					}

				case _, ok := <-cluster.rpcServicePortsDiscoveryChanResp:
					if ok {
						cluster.bootstrapExpectedSize.Add(1)
					}
				}
			}
		})

		go cluster.checkBootstrapSize()
		go func() {
			time.Sleep(2 * time.Second)
			cancel()
		}()
		cluster.wg.Wait()
	})

	t.Run("check_bootstrap_size_timer_retry", func(t *testing.T) {
		clusters := makeSizedCluster(sizedClusterConfig{})
		cluster := clusters[0]
		ctx, cancel := context.WithCancel(context.Background())
		cluster.ctx = ctx
		cluster.checkBootstrapSizeDuration = 100 * time.Millisecond
		cluster.di.sendRPCFunc = func(address string, client scalezillapb.ScalezillaClient, request RPCRequest) {
			time.Sleep(600 * time.Millisecond)
			cluster.rpcServicePortsDiscoveryChanResp <- RPCResponse{}
		}

		cluster.wg.Go(func() {
			for {
				select {
				case <-cluster.ctx.Done():
					return

				case data, ok := <-cluster.rpcServicePortsDiscoveryChanReq:
					if ok {
						cluster.rcvServicePortsDiscovery(data)
					}

				case _, ok := <-cluster.rpcServicePortsDiscoveryChanResp:
					if ok {
						fmt.Println("receive rpc")
						cluster.bootstrapExpectedSize.Add(1)
					}
				}
			}
		})

		go cluster.checkBootstrapSize()
		go func() {
			time.Sleep(2 * time.Second)
			cancel()
		}()
		cluster.wg.Wait()
	})
}

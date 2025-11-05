package cluster

import (
	"errors"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc(t *testing.T) {
	assert := assert.New(t)

	t.Run("start_grpc_server_listen_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		cluster.grpcListenFunc = func(network, address string) (net.Listener, error) {
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

		cluster.grpcListenFunc = func(network, address string) (net.Listener, error) {
			return badListener{}, nil
		}
		cluster.grpcAddress = net.TCPAddr{}

		assert.Error(cluster.startGRPCServer())
	})

	t.Run("start_stop_grpc_server_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.buildAddressAndID()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		go func() {
			assert.Nil(cluster.startGRPCServer())
		}()

		time.Sleep(100 * time.Millisecond)
		cluster.stopGRPCServer()
	})
}

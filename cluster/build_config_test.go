package cluster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_build_config(t *testing.T) {
	assert := assert.New(t)

	t.Run("build_address_and_id", func(t *testing.T) {
		cluster := &Cluster{
			config: Config{
				HostIPAddress: defaultHostIPAddress,
				RaftGRPCPort:  defaultRaftGRPCPort,
			},
		}

		cluster.buildAddressAndID()
		assert.NotNil(cluster.grpcAddress.String())
		assert.NotNil(cluster.raftyAddress.String())
		assert.NotNil(cluster.id)
	})

	t.Run("build_peers", func(t *testing.T) {
		cluster := &Cluster{
			config: Config{
				HostIPAddress: defaultHostIPAddress,
				RaftGRPCPort:  defaultRaftGRPCPort,
			},
			nodeMap: make(map[string]*nodeMap),
		}
		httpPort := 16000
		grpcPort := 16001
		raftPort := 16002
		address := "127.0.0.1"
		cluster.members = append(cluster.members, fmt.Sprintf("%s:%d", address, grpcPort))
		cluster.nodeMap["16000"] = &nodeMap{
			IsVoter:   true,
			ID:        fmt.Sprintf("%d", grpcPort),
			Address:   address,
			HTTPPort:  uint32(httpPort),
			GRPCPort:  uint32(grpcPort),
			RaftyPort: uint32(raftPort),
			NodePool:  defaultNodePool,
		}

		assert.NotNil(cluster.buildPeers())
	})

	t.Run("build_datadir", func(t *testing.T) {
		cluster := &Cluster{}

		cluster.buildDataDir()
		assert.NotNil(cluster.config.DataDir)
	})

	t.Run("build_store", func(t *testing.T) {
		cluster := &Cluster{}

		cluster.buildDataDir()
		store, err := cluster.buildStore()
		assert.Nil(err)
		assert.Nil(store.Close())
	})

	t.Run("build_dev_config", func(t *testing.T) {
		config := ClusterInitialConfig{Dev: true}

		cluster, err := NewCluster(config)
		assert.Nil(err)
		assert.Equal(cluster.dev, true)
	})
}

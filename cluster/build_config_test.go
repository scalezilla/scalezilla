package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildConfig(t *testing.T) {
	assert := assert.New(t)

	t.Run("build_address_and_id", func(t *testing.T) {
		cluster := &Cluster{
			config: Config{
				HostIPAddress: defaultHostIPAddress,
				RaftGRPCPort:  defaultRaftGRPCPort,
			},
		}

		cluster.buildAddressAndID()
		assert.NotNil(cluster.address.String())
		assert.NotNil(cluster.id)
	})

	t.Run("build_peers", func(t *testing.T) {
		cluster := &Cluster{
			config: Config{
				HostIPAddress: defaultHostIPAddress,
				RaftGRPCPort:  defaultRaftGRPCPort,
			},
		}
		cluster.members = append(cluster.members, "127.0.0.1:16000")

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

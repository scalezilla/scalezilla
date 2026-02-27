package cluster

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_utils(t *testing.T) {
	assert := assert.New(t)

	t.Run("get_service_address_from_raft", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		_, ok := cluster.getServiceAddressFromRaft("")
		assert.Equal(false, ok)

		_, ok = cluster.getServiceAddressFromRaft("10.0.0.1:50000")
		assert.Equal(false, ok)

		cluster.nodeMap[cluster.members_grpc[0]] = &nodeMap{
			HTTPPort:  uint32(defaultHTTPPort),
			GRPCPort:  uint32(defaultGRPCPort),
			RaftyPort: uint32(defaultRaftGRPCPort),
		}
		_, ok = cluster.getServiceAddressFromRaft(cluster.members_raft[0])
		assert.Equal(true, ok)
	})
}

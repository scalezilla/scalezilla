package cluster

import (
	"os"
	"testing"

	"github.com/scalezilla/scalezilla/osdiscovery"
	"github.com/stretchr/testify/assert"
)

func TestCluster_system_info(t *testing.T) {
	assert := assert.New(t)

	t.Run("check_system_info_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		assert.Nil(cluster.checkSystemInfo())
	})

	t.Run("check_system_info_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		cluster.osdiscoveryFunc = func() *osdiscovery.SystemInfo {
			return &osdiscovery.SystemInfo{
				Cgroups: &osdiscovery.Cgroups{
					Version: 1,
				},
			}
		}

		assert.Error(cluster.checkSystemInfo())
	})
}

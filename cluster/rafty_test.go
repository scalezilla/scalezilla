package cluster

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCluster_rafty(t *testing.T) {
	assert := assert.New(t)

	t.Run("new_rafty_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		dir := filepath.Dir(cluster.config.DataDir)
		assert.Nil(os.MkdirAll(dir, 0750))
		file, err := os.Create(cluster.config.DataDir)
		assert.Nil(err)
		assert.Nil(file.Close())
		cluster.raftMetricPrefix = "new_rafty_error"
		_, err = cluster.newRafty()
		assert.Error(err)
	})

	t.Run("new_rafty_success_voter", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		cluster.raftMetricPrefix = "new_rafty_success_voter"
		_, err := cluster.newRafty()
		assert.Nil(err)
	})

	t.Run("new_rafty_success_non_voter", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		cluster.isVoter = false
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		cluster.raftMetricPrefix = "new_rafty_success_non_voter"
		_, err := cluster.newRafty()
		assert.Nil(err)
	})

	t.Run("start_failure", func(t *testing.T) {
		mock := mockRafty{
			err: errors.New("start error"),
		}
		cluster := &Cluster{
			rafty: &mock,
		}

		assert.Error(cluster.startRafty())
		assert.Equal(true, mock.called)
	})

	t.Run("stop_success", func(t *testing.T) {
		mock := mockRafty{}
		cluster := &Cluster{
			rafty: &mock,
		}

		cluster.stopRafty()
		assert.Equal(true, mock.called)
	})

	t.Run("start_rafty_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		var err error
		cluster.raftMetricPrefix = "start_rafty_success"
		cluster.rafty, err = cluster.newRafty()
		assert.Nil(err)

		defer func() {
			time.Sleep(200 * time.Millisecond)
			cluster.stopRafty()
		}()
		assert.Nil(cluster.startRafty())
	})
}

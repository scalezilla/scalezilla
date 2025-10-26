package cluster

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRafty(t *testing.T) {
	assert := assert.New(t)

	t.Run("new_rafty_error", func(t *testing.T) {
		cluster := makeBasicCluster(false)
		defer func() {
			_ = os.RemoveAll(cluster.dataDir)
		}()

		dir := filepath.Dir(cluster.dataDir)
		assert.Nil(os.MkdirAll(dir, 0750))
		file, err := os.Create(cluster.dataDir)
		assert.Nil(err)
		assert.Nil(file.Close())
		_, err = cluster.newRafty("new_rafty_error")
		assert.Error(err)
	})

	t.Run("new_rafty_success", func(t *testing.T) {
		cluster := makeBasicCluster(false)
		defer func() {
			_ = os.RemoveAll(cluster.dataDir)
		}()
		_, err := cluster.newRafty("new_rafty_success")
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
		cluster := makeBasicCluster(true)
		defer func() {
			_ = os.RemoveAll(cluster.dataDir)
		}()

		var err error
		cluster.rafty, err = cluster.newRafty("start_rafty_success")
		assert.Nil(err)

		defer func() {
			time.Sleep(200 * time.Millisecond)
			cluster.stopRafty()
		}()
		assert.Nil(cluster.startRafty())
	})
}

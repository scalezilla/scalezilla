package cluster

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClusterStart(t *testing.T) {
	assert := assert.New(t)

	t.Run("start_new_rafty_error", func(t *testing.T) {
		cluster := makeBasicCluster(false)
		defer func() {
			_ = os.RemoveAll(cluster.dataDir)
		}()

		dir := filepath.Dir(cluster.dataDir)
		assert.Nil(os.MkdirAll(dir, 0750))
		file, err := os.Create(cluster.dataDir)
		assert.Nil(err)
		assert.Nil(file.Close())
		assert.Error(cluster.Start())
	})

	t.Run("start_rafty_func_error", func(t *testing.T) {
		cluster := makeBasicCluster(false)
		defer func() {
			assert.Nil(cluster.raftyStore.Close())
			_ = os.RemoveAll(cluster.dataDir)
		}()
		cluster.startRaftyFunc = func() error {
			return errors.New("start error")
		}
		assert.Error(cluster.Start())
	})

	t.Run("start_api_server_func_error", func(t *testing.T) {
		cluster := makeBasicCluster(false)
		defer func() {
			assert.Nil(cluster.raftyStore.Close())
			_ = os.RemoveAll(cluster.dataDir)
		}()
		cluster.startRaftyFunc = func() error {
			return nil
		}

		cluster.startAPIServerFunc = func() error {
			return errors.New("start error")
		}
		assert.Error(cluster.Start())
	})

	t.Run("stop_api_server_func_error", func(t *testing.T) {
		cluster := makeBasicCluster(false)
		defer func() {
			assert.Nil(cluster.raftyStore.Close())
			_ = os.RemoveAll(cluster.dataDir)
		}()

		cluster.startRaftyFunc = func() error {
			return nil
		}
		cluster.startAPIServerFunc = func() error {
			return nil
		}

		cluster.stopAPIServerFunc = func() error {
			return errors.New("start error")
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			close(cluster.quit)
		}()
		assert.Error(cluster.Start())
	})

	t.Run("start_success", func(t *testing.T) {
		cluster := makeBasicCluster(false)
		defer func() {
			assert.Nil(cluster.raftyStore.Close())
			_ = os.RemoveAll(cluster.dataDir)
		}()

		cluster.startRaftyFunc = func() error {
			return nil
		}
		cluster.startAPIServerFunc = func() error {
			return nil
		}

		cluster.stopAPIServerFunc = func() error {
			return nil
		}

		cluster.stopRaftyFunc = func() {}

		go func() {
			time.Sleep(100 * time.Millisecond)
			close(cluster.quit)
		}()
		assert.Nil(cluster.Start())
	})

	t.Run("rafty_store_close_error", func(t *testing.T) {
		cluster := makeBasicCluster(false)
		defer func() {
			assert.Nil(cluster.raftyStore.Close())
			_ = os.RemoveAll(cluster.dataDir)
		}()

		cluster.startRaftyFunc = func() error {
			return nil
		}
		cluster.startAPIServerFunc = func() error {
			return nil
		}

		cluster.stopAPIServerFunc = func() error {
			return nil
		}

		cluster.stopRaftyFunc = func() {}

		cluster.raftyStoreCloseFunc = func() error {
			return errors.New("close store error")
		}

		go func() {
			time.Sleep(100 * time.Millisecond)
			close(cluster.quit)
		}()
		assert.Error(cluster.Start())
	})
}

package cluster

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Lord-Y/rafty"
	"github.com/stretchr/testify/assert"
)

func TestCluster(t *testing.T) {
	assert := assert.New(t)

	t.Run("start_new_rafty_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		cluster.newRaftyFunc = func() (*rafty.Rafty, error) {
			return nil, errors.New("start error")
		}
		assert.Error(cluster.Start())
	})

	t.Run("start_rafty_func_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		cluster.startRaftyFunc = func() error {
			return errors.New("start error")
		}
		assert.Error(cluster.Start())
	})

	t.Run("start_api_server_func_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
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
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		sigCtx, stop := BuildSignal(context.Background())
		cluster.ctx = sigCtx
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
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
			assert.Error(cluster.Start())
		}()

		time.Sleep(300 * time.Millisecond)
		stop()
	})

	t.Run("start_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		sigCtx, stop := BuildSignal(context.Background())
		cluster.ctx = sigCtx
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
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
			assert.Nil(cluster.Start())
		}()

		time.Sleep(300 * time.Millisecond)
		stop()
	})

	t.Run("rafty_store_close_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		sigCtx, stop := BuildSignal(context.Background())
		cluster.ctx = sigCtx
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
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
			assert.Error(cluster.Start())
		}()

		time.Sleep(300 * time.Millisecond)
		stop()
	})
}

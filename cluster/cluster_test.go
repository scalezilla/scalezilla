package cluster

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Lord-Y/rafty"
	"github.com/stretchr/testify/assert"
)

func TestCluster(t *testing.T) {
	assert := assert.New(t)

	t.Run("new_cluster_config_success", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		configDir := "testdata/config"
		configFile := filepath.Join(workingDir, configDir, "config_success_server.hcl")
		config := ClusterInitialConfig{
			ConfigFile: configFile,
			Test:       true,
		}

		_, err = NewCluster(config)
		assert.Nil(err)
	})

	t.Run("new_cluster_config_error", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		configDir := "testdata/config"
		configFile := filepath.Join(workingDir, configDir, "config_error_empty.hcl")
		config := ClusterInitialConfig{
			ConfigFile: configFile,
			Test:       true,
		}

		_, err = NewCluster(config)
		assert.Error(err)
	})

	t.Run("start_system_info_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		cluster.checkSystemInfoFunc = func() error {
			return errors.New("start error")
		}
		assert.Error(cluster.Start())
	})

	t.Run("start_new_rafty_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		cluster.checkSystemInfoFunc = func() error {
			return nil
		}

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

		cluster.checkSystemInfoFunc = func() error {
			return nil
		}

		cluster.newRaftyFunc = func() (*rafty.Rafty, error) {
			return nil, nil
		}

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

		cluster.checkSystemInfoFunc = func() error {
			return nil
		}

		cluster.newRaftyFunc = func() (*rafty.Rafty, error) {
			return nil, nil
		}

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

		cluster.checkSystemInfoFunc = func() error {
			return nil
		}

		cluster.newRaftyFunc = func() (*rafty.Rafty, error) {
			return nil, nil
		}

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

		cluster.checkSystemInfoFunc = func() error {
			return nil
		}

		// NO dependency injection here
		// cluster.newRaftyFunc = func() (*rafty.Rafty, error) {
		// 	return nil, nil
		// }

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

	t.Run("start_rafty_store_close_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		sigCtx, stop := BuildSignal(context.Background())
		cluster.ctx = sigCtx
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		cluster.checkSystemInfoFunc = func() error {
			return nil
		}

		cluster.newRaftyFunc = func() (*rafty.Rafty, error) {
			return nil, nil
		}

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

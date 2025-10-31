package commands

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandsConfig(t *testing.T) {
	assert := assert.New(t)

	t.Run("start_success", func(t *testing.T) {
		configDir := "../cluster/testdata/config"
		configFile := filepath.Join(configDir, "config_success_server.hcl")

		cmd := Config()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run(ctx, []string{"config", "--file", configFile, "--test-raft-metric-prefix", "commands-config-success", "--test"})
		}()

		time.Sleep(time.Second)
		cancel()

		select {
		case err := <-done:
			assert.NoError(err)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for Run() to stop")
		}
	})

	t.Run("start_error", func(t *testing.T) {
		configDir := "../cluster/testdata/config"
		configFile := filepath.Join(configDir, "config_error_empty.hcl")

		cmd := Config()
		assert.Error(cmd.Run(context.Background(), []string{"config", "--file", configFile, "--test-raft-metric-prefix", "commands-config-error"}))
	})

	t.Run("validate", func(t *testing.T) {
		configDir := "../cluster/testdata/config"
		configFile := filepath.Join(configDir, "config_success_server.hcl")

		cmd := Config()

		done := make(chan error, 1)
		done <- cmd.Run(context.Background(), []string{"config", "--validate", "--file", configFile, "--test-raft-metric-prefix", "commands-config-success", "--test"})

	})
}

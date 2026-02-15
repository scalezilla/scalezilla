package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandsNodesList(t *testing.T) {
	assert := assert.New(t)

	t.Run("success", func(t *testing.T) {
		cmd := Nodes()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run(ctx, []string{"nodes", "list"})
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

	t.Run("error", func(t *testing.T) {
		cmd := Nodes()
		assert.Error(cmd.Run(context.Background(), []string{"nodes", "list", "--output", "zzz"}))
	})
}

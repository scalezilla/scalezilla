package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandsDev(t *testing.T) {
	assert := assert.New(t)

	t.Run("success", func(t *testing.T) {
		cmd := Dev()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run(ctx, []string{"dev", "--test-raft-metric-prefix", "commands-dev-success"})
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

	t.Run("fail", func(t *testing.T) {
		cmd := Dev()
		assert.Error(cmd.Run(context.Background(), []string{"dev", "--fail"}))
	})
}

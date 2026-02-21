package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandsBootstrap(t *testing.T) {
	assert := assert.New(t)

	t.Run("status", func(t *testing.T) {
		cmd := Bootstrap()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run(ctx, []string{"bootstrap", "status"})
		}()

		time.Sleep(time.Second)
		cancel()

		select {
		case err := <-done:
			assert.Error(err)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for Run() to stop")
		}
	})

	t.Run("bootatrap", func(t *testing.T) {
		cmd := Bootstrap()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run(ctx, []string{"bootstrap", "cluster"})
		}()

		time.Sleep(time.Second)
		cancel()

		select {
		case err := <-done:
			assert.Error(err)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for Run() to stop")
		}
	})
}

package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandsPods(t *testing.T) {
	assert := assert.New(t)

	t.Run("list", func(t *testing.T) {
		cmd := Pods()
		assert.Error(cmd.Run(context.Background(), []string{"pods", "list"}))
	})

	t.Run("list_output", func(t *testing.T) {
		cmd := Pods()
		assert.Error(cmd.Run(context.Background(), []string{"pods", "list", "--output", "zzz"}))
	})
}

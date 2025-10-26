package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandsAgent(t *testing.T) {
	assert := assert.New(t)

	t.Run("help", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer t.Cleanup(cancel)

		cmd := Agent()
		assert.Nil(cmd.Run(ctx, []string{"-h"}))
	})
}

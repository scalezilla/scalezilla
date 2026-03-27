package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandsDeployment(t *testing.T) {
	assert := assert.New(t)

	t.Run("apply", func(t *testing.T) {
		cmd := Deployment()
		assert.Error(cmd.Run(context.Background(), []string{"deployment", "apply", "--file", "zzz"}))
	})
}

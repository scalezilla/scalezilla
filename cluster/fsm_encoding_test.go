package cluster

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_fsm_encoding(t *testing.T) {
	assert := assert.New(t)

	t.Run("decode_command_success", func(t *testing.T) {
		cmd := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
			Token:      "b",
		}
		buffer := new(bytes.Buffer)
		err := aclTokenEncodeCommand(cmd, buffer)
		assert.Nil(err)
		z, err := decodeCommand(buffer.Bytes())
		assert.Equal(cmd.Kind, z)
	})

	t.Run("decode_command_error", func(t *testing.T) {
		_, err := decodeCommand(nil)
		assert.Error(err)
	})
}

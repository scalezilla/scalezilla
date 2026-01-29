package cluster

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLuster_api_submit_commands(t *testing.T) {
	assert := assert.New(t)

	t.Run("submit_command_acl_token_write_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)
		token := &AclToken{
			AccessorID: "a",
			Token:      "b",
		}
		mock := mockRafty{}
		cluster.rafty = &mock

		assert.Nil(cluster.submitCommandACLTokenWrite(aclTokenCommandSet, token))
		assert.Equal(true, mock.called)
	})

	t.Run("submit_command_acl_token_write_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)
		token := &AclToken{
			AccessorID: "a",
			Token:      "b",
		}

		cluster.di.aclTokenEncodeCommandFunc = func(cmd aclTokenCommand, w io.Writer) error {
			return errors.New("acl encode commnd error")
		}
		assert.Error(cluster.submitCommandACLTokenWrite(aclTokenCommandSet, token))

		mock := mockRafty{
			err: errors.New("rafty submit command error"),
		}
		cluster.rafty = &mock

		cluster.di.aclTokenEncodeCommandFunc = aclTokenEncodeCommand

		assert.Error(cluster.submitCommandACLTokenWrite(aclTokenCommandSet, token))
		assert.Equal(true, mock.called)
	})
}

package cluster

import (
	"bytes"
	"os"
	"testing"

	"github.com/Lord-Y/rafty"
	"github.com/stretchr/testify/assert"
)

func TestCluster_fsm_acl_token(t *testing.T) {
	assert := assert.New(t)

	t.Run("acl_token_set", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		cmd := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
			Token:      "b",
		}
		buffer := new(bytes.Buffer)
		err = aclTokenEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		assert.Nil(cluster.fsm.memoryStore.aclTokenSet(entry, cmd))
		assert.Nil(cluster.fsm.memoryStore.aclTokenSet(entry, cmd)) // for delete
	})

	t.Run("acl_token_get", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		cmd := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
			Token:      "b",
		}
		buffer := new(bytes.Buffer)
		err = aclTokenEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		r, err := cluster.fsm.memoryStore.aclTokenGet([]byte("a"))
		assert.Nil(r)
		assert.ErrorIs(err, rafty.ErrKeyNotFound)
		assert.Nil(cluster.fsm.memoryStore.aclTokenSet(entry, cmd))
		r, err = cluster.fsm.memoryStore.aclTokenGet([]byte("a"))
		assert.NotNil(r)
		assert.Nil(err)
	})

	t.Run("acl_token_exist", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		cmd := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
			Token:      "b",
		}
		buffer := new(bytes.Buffer)
		err = aclTokenEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		assert.Equal(false, cluster.fsm.memoryStore.aclTokenExist([]byte("a")))
		assert.Nil(cluster.fsm.memoryStore.aclTokenSet(entry, cmd))
		assert.Equal(true, cluster.fsm.memoryStore.aclTokenExist([]byte("a")))
	})

	t.Run("acl_token_delete", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		cmd := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
			Token:      "b",
		}
		buffer := new(bytes.Buffer)
		err = aclTokenEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		cluster.fsm.memoryStore.aclTokenDelete([]byte("a"))
		assert.Nil(cluster.fsm.memoryStore.aclTokenSet(entry, cmd))
		cluster.fsm.memoryStore.aclTokenDelete([]byte("a"))
	})

	t.Run("acl_token_get_all", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		cmd := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
			Token:      "b",
		}
		buffer := new(bytes.Buffer)
		err = aclTokenEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		r, err := cluster.fsm.memoryStore.aclTokenGetAll()
		assert.Nil(r)
		assert.Nil(err)
		assert.Nil(cluster.fsm.memoryStore.aclTokenSet(entry, cmd))
		r, err = cluster.fsm.memoryStore.aclTokenGetAll()
		assert.NotNil(r)
		assert.Nil(err)
	})

	t.Run("acl_token_encoded", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		set := aclTokenCommand{
			Kind:       aclTokenCommandSet,
			AccessorID: "a",
			Token:      "b",
		}
		get := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
		}
		fakeGet := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "aa",
		}
		getAll := aclTokenCommand{
			Kind:       aclTokenCommandGetAll,
			AccessorID: "a",
		}
		buffer := new(bytes.Buffer)
		err = aclTokenEncodeCommand(set, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Index:   1,
			Command: buffer.Bytes(),
		}

		assert.Nil(cluster.fsm.memoryStore.aclTokenSet(entry, set))
		r, err := cluster.fsm.memoryStore.aclTokenEncoded(get)
		assert.NotNil(r)
		assert.Nil(err)

		r, err = cluster.fsm.memoryStore.aclTokenEncoded(fakeGet)
		assert.Nil(r)
		assert.NotNil(err, rafty.ErrKeyNotFound)

		r, err = cluster.fsm.memoryStore.aclTokenEncoded(getAll)
		assert.NotNil(r)
		assert.Nil(err)
	})
}

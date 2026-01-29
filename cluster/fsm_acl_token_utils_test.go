package cluster

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"testing"

	"github.com/Lord-Y/rafty"
	"github.com/stretchr/testify/assert"
)

func TestCluster_fsm_acl_token_utils(t *testing.T) {
	assert := assert.New(t)

	w := &failWriter{failOn: 1}
	buffer := new(bytes.Buffer)

	t.Run("acl_token_encode_command", func(t *testing.T) {
		cmd := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
			Token:      "b",
		}
		assert.Error(aclTokenEncodeCommand(cmd, w))

		// write binary accessor id
		w = &failWriter{failOn: 2}
		assert.Error(aclTokenEncodeCommand(cmd, w))

		// write cmd accessor id
		w = &failWriter{failOn: 3}
		assert.Error(aclTokenEncodeCommand(cmd, w))

		// write binary token
		w = &failWriter{failOn: 4}
		assert.Error(aclTokenEncodeCommand(cmd, w))

		// write cmd token
		w = &failWriter{failOn: 5}
		assert.Error(aclTokenEncodeCommand(cmd, w))

		// write binary initial token
		w = &failWriter{failOn: 5}
		assert.Error(aclTokenEncodeCommand(cmd, w))

		// write initial token
		w = &failWriter{failOn: 6}
		assert.Error(aclTokenEncodeCommand(cmd, w))

		// No errors expected here
		assert.Nil(aclTokenEncodeCommand(cmd, buffer))
		assert.NotNil(buffer.Bytes())
	})

	t.Run("acl_token_decode_command", func(t *testing.T) {
		cmd := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
			Token:      "b",
		}
		// error kind
		_, err := aclTokenDecodeCommand([]byte{})
		assert.Error(err)

		// error accessor id length
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // kind
		_, err = aclTokenDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error accessor id
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(1)) // accessor id len
		_, err = aclTokenDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error token length
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(1)) // accessor id len
		_, _ = buf.Write([]byte("a"))                         // accessor id
		_, err = aclTokenDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error token
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(1)) // accessor id len
		_, _ = buf.Write([]byte("a"))                         // accessor id
		_ = binary.Write(buf, binary.LittleEndian, uint64(1)) // token len
		_, err = aclTokenDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error initial token
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(1)) // accessor id len
		_, _ = buf.Write([]byte("a"))                         // accessor id
		_ = binary.Write(buf, binary.LittleEndian, uint64(1)) // token len
		_, _ = buf.Write([]byte("b"))                         // token
		_, err = aclTokenDecodeCommand(buf.Bytes())
		assert.Error(err)

		// No errors expected here
		dec, err := aclTokenDecodeCommand(buffer.Bytes())
		assert.Nil(err)
		assert.Equal(cmd, dec)
	})

	t.Run("acl_token_apply_command_linearizable_read", func(t *testing.T) {
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
			LogType: uint32(rafty.LogCommandLinearizableRead),
			Term:    1,
			Command: buffer.Bytes(),
		}

		r, err := cluster.fsm.aclTokenApplyCommand(entry)
		assert.Nil(r)
		assert.Nil(err)
	})

	t.Run("acl_token_apply_command_set", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		cmd := aclTokenCommand{
			Kind:       aclTokenCommandSet,
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

		r, err := cluster.fsm.aclTokenApplyCommand(entry)
		assert.Nil(r)
		assert.Nil(err)
	})

	t.Run("acl_token_apply_command_get", func(t *testing.T) {
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

		r, err := cluster.fsm.aclTokenApplyCommand(entry)
		assert.Nil(r)
		assert.ErrorIs(err, rafty.ErrKeyNotFound)
	})

	t.Run("acl_token_apply_command_all", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		cmds := aclTokenCommand{
			Kind:       aclTokenCommandSet,
			AccessorID: "a",
			Token:      "b",
		}

		buffers := new(bytes.Buffer)
		err = aclTokenEncodeCommand(cmds, buffers)
		assert.Nil(err)
		entrySet := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffers.Bytes(),
		}

		r, err := cluster.fsm.aclTokenApplyCommand(entrySet)
		assert.Nil(r)
		assert.Nil(err)

		cmdg := aclTokenCommand{
			Kind:       aclTokenCommandGet,
			AccessorID: "a",
		}

		bufferg := new(bytes.Buffer)
		err = aclTokenEncodeCommand(cmdg, bufferg)
		assert.Nil(err)
		entryGet := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: bufferg.Bytes(),
		}

		r, err = cluster.fsm.aclTokenApplyCommand(entryGet)
		assert.NotNil(r)
		assert.Nil(err)

		cmdd := aclTokenCommand{
			Kind:       aclTokenCommandDelete,
			AccessorID: "a",
		}

		bufferd := new(bytes.Buffer)
		err = aclTokenEncodeCommand(cmdd, bufferd)
		assert.Nil(err)
		entryDel := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: bufferd.Bytes(),
		}

		r, err = cluster.fsm.aclTokenApplyCommand(entryDel)
		assert.Nil(r)
		assert.Nil(err)
	})

	t.Run("acl_token_unmarshal_success", func(t *testing.T) {
		z := AclToken{
			AccessorID: "a",
			Token:      "a",
		}
		r, err := json.Marshal(z)
		assert.Nil(err)
		ra, err := aclTokenUnmarshal(r)
		assert.Equal(z, ra)
		assert.Nil(err)
	})

	t.Run("acl_token_unmarshal_error", func(t *testing.T) {
		_, err := aclTokenUnmarshal([]byte("a"))
		assert.Error(err)
	})
}

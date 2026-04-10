package cluster

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Lord-Y/rafty"
	"github.com/stretchr/testify/assert"
)

func TestCluster_fsm_deployment_utils(t *testing.T) {
	assert := assert.New(t)

	w := &failWriter{failOn: 1}
	buffer := new(bytes.Buffer)

	t.Run("encode_command", func(t *testing.T) {
		content := make(map[uint64]deploymentContent)
		content[1] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now(),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind:               deploymentCommandSet,
			Name:               "redis",
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write binary name
		w = &failWriter{failOn: 2}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write cmd name
		w = &failWriter{failOn: 3}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write newRollingVersion
		w = &failWriter{failOn: 4}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write currentUsedVersion
		w = &failWriter{failOn: 5}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write binary count
		w = &failWriter{failOn: 5}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write key
		w = &failWriter{failOn: 6}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write isStable
		w = &failWriter{failOn: 7}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write isStable
		w = &failWriter{failOn: 8}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write rawContent
		w = &failWriter{failOn: 9}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write rawContent
		w = &failWriter{failOn: 10}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write version
		w = &failWriter{failOn: 11}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write createdAt
		w = &failWriter{failOn: 12}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write replicaSetID
		w = &failWriter{failOn: 13}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write replicaSetID
		w = &failWriter{failOn: 14}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// write mustBeStarted
		w = &failWriter{failOn: 15}
		assert.Error(deploymentEncodeCommand(cmd, w))

		// No errors expected here
		assert.Nil(deploymentEncodeCommand(cmd, buffer))
		assert.NotNil(buffer.Bytes())
	})

	t.Run("decode_command", func(t *testing.T) {
		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind:               deploymentCommandSet,
			Name:               "redis",
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		// error kind
		_, err := deploymentDecodeCommand([]byte{})
		assert.Error(err)

		// error name length
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind)) // kind
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error name
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))      // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name))) // name len
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error newRollingVersion
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))      // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name))) // name len
		_, _ = buf.Write([]byte(cmd.Name))                                // name
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error currentUsedVersion
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))      // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name))) // name len
		_, _ = buf.Write([]byte(cmd.Name))                                // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion) // newRollingVersion
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error count
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))       // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))  // name len
		_, _ = buf.Write([]byte(cmd.Name))                                 // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)  // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion) // currentUsedVersion
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error key
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))         // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))    // name len
		_, _ = buf.Write([]byte(cmd.Name))                                   // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)    // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)   // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content))) // count
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error isStable
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))         // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))    // name len
		_, _ = buf.Write([]byte(cmd.Name))                                   // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)    // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)   // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content))) // count
		_ = binary.Write(buf, binary.LittleEndian, key)                      // key
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error rawContentLen
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))         // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))    // name len
		_, _ = buf.Write([]byte(cmd.Name))                                   // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)    // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)   // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content))) // count
		_ = binary.Write(buf, binary.LittleEndian, key)                      // key
		_ = binary.Write(buf, binary.LittleEndian, uint64(1))                // isStable
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error rawContent
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))                                 // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))                            // name len
		_, _ = buf.Write([]byte(cmd.Name))                                                           // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)                            // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)                           // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content)))                         // count
		_ = binary.Write(buf, binary.LittleEndian, key)                                              // key
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].IsStable)                // isStable
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].RawContent))) // rawContentLen len
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error version
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))                                 // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))                            // name len
		_, _ = buf.Write([]byte(cmd.Name))                                                           // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)                            // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)                           // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content)))                         // count
		_ = binary.Write(buf, binary.LittleEndian, key)                                              // key
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].IsStable)                // isStable
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].RawContent))) // rawContentLen len
		_, _ = buf.Write([]byte(cmd.Content[uint64(key)].RawContent))                                // rawContent
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error createdAt
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))                                 // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))                            // name len
		_, _ = buf.Write([]byte(cmd.Name))                                                           // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)                            // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)                           // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content)))                         // count
		_ = binary.Write(buf, binary.LittleEndian, key)                                              // key
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].IsStable)                // isStable
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].RawContent))) // rawContentLen len
		_, _ = buf.Write([]byte(cmd.Content[uint64(key)].RawContent))                                // rawContent
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].Version)                 // version
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error replicaSetIDLen
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))                                 // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))                            // name len
		_, _ = buf.Write([]byte(cmd.Name))                                                           // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)                            // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)                           // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content)))                         // count
		_ = binary.Write(buf, binary.LittleEndian, key)                                              // key
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].IsStable)                // isStable
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].RawContent))) // rawContentLen len
		_, _ = buf.Write([]byte(cmd.Content[uint64(key)].RawContent))                                // rawContent
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].Version)                 // version
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].CreatedAt.UnixNano())    // createdAt
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error replicaSetID
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))                                   // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))                              // name len
		_, _ = buf.Write([]byte(cmd.Name))                                                             // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)                              // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)                             // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content)))                           // count
		_ = binary.Write(buf, binary.LittleEndian, key)                                                // key
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].IsStable)                  // isStable
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].RawContent)))   // rawContentLen len
		_, _ = buf.Write([]byte(cmd.Content[uint64(key)].RawContent))                                  // rawContent
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].Version)                   // version
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].CreatedAt.UnixNano())      // createdAt
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].ReplicaSetID))) // replicaSetID len
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// error mustBeStarted
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))                                   // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))                              // name len
		_, _ = buf.Write([]byte(cmd.Name))                                                             // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)                              // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)                             // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content)))                           // count
		_ = binary.Write(buf, binary.LittleEndian, key)                                                // key
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].IsStable)                  // isStable
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].RawContent)))   // rawContentLen len
		_, _ = buf.Write([]byte(cmd.Content[uint64(key)].RawContent))                                  // rawContent
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].Version)                   // version
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].CreatedAt.UnixNano())      // createdAt
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].ReplicaSetID))) // replicaSetID len
		_, _ = buf.Write([]byte(cmd.Content[uint64(key)].ReplicaSetID))                                // replicaSetID
		_, err = deploymentDecodeCommand(buf.Bytes())
		assert.Error(err)

		// No errors expected here
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint32(cmd.Kind))                                   // kind
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Name)))                              // name len
		_, _ = buf.Write([]byte(cmd.Name))                                                             // name
		_ = binary.Write(buf, binary.LittleEndian, cmd.NewRollingVersion)                              // newRollingVersion
		_ = binary.Write(buf, binary.LittleEndian, cmd.CurrentUsedVersion)                             // currentUsedVersion
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content)))                           // count
		_ = binary.Write(buf, binary.LittleEndian, key)                                                // key
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].IsStable)                  // isStable
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].RawContent)))   // rawContentLen len
		_, _ = buf.Write([]byte(cmd.Content[uint64(key)].RawContent))                                  // rawContent
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].Version)                   // version
		_ = binary.Write(buf, binary.LittleEndian, cmd.Content[uint64(key)].CreatedAt.UnixNano())      // createdAt
		_ = binary.Write(buf, binary.LittleEndian, uint64(len(cmd.Content[uint64(key)].ReplicaSetID))) // replicaSetID len
		_, _ = buf.Write([]byte(cmd.Content[uint64(key)].ReplicaSetID))                                // replicaSetID
		_ = binary.Write(buf, binary.LittleEndian, cmd.MustBeStarted)                                  // mustBeStarted
		dec, err := deploymentDecodeCommand(buf.Bytes())
		assert.NoError(err)
		assert.Equal(cmd, dec)
	})

	t.Run("apply_command_linearizable_read", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind: deploymentCommandGet,
			Name: "redis",
		}

		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogCommandLinearizableRead),
			Term:    1,
			Command: buffer.Bytes(),
		}

		r, err := cluster.fsm.deploymentApplyCommand(entry)
		assert.Nil(r)
		assert.Nil(err)
	})

	t.Run("apply_command_set", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind:               deploymentCommandSet,
			Name:               "redis",
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		r, err := cluster.fsm.deploymentApplyCommand(entry)
		assert.Nil(r)
		assert.Nil(err)
	})

	t.Run("apply_command_get", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind:               deploymentCommandGet,
			Name:               "redis",
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		r, err := cluster.fsm.deploymentApplyCommand(entry)
		assert.Nil(r)
		assert.ErrorIs(err, rafty.ErrKeyNotFound)
	})

	t.Run("apply_command_all", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmds := deploymentState{
			Kind:               deploymentCommandSet,
			Name:               "redis",
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		buffers := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmds, buffers)
		assert.Nil(err)
		entrySet := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffers.Bytes(),
		}

		r, err := cluster.fsm.deploymentApplyCommand(entrySet)
		assert.Nil(r)
		assert.Nil(err)

		cmdg := deploymentState{
			Kind: deploymentCommandGet,
			Name: "redis",
		}

		bufferg := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmdg, bufferg)
		assert.Nil(err)
		entryGet := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: bufferg.Bytes(),
		}

		r, err = cluster.fsm.deploymentApplyCommand(entryGet)
		assert.NotNil(r)
		assert.Nil(err)

		cmdd := deploymentState{
			Kind: deploymentCommandDelete,
			Name: "redis",
		}

		bufferd := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmdd, bufferd)
		assert.Nil(err)
		entryDel := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: bufferd.Bytes(),
		}

		r, err = cluster.fsm.deploymentApplyCommand(entryDel)
		assert.Nil(r)
		assert.Nil(err)
	})

	t.Run("unmarshal_success", func(t *testing.T) {
		z := deploymentState{
			Name: "redis",
		}
		r, err := json.Marshal(z)
		assert.Nil(err)
		ra, err := deploymentUnmarshal(r)
		assert.Equal(z, ra)
		assert.Nil(err)
	})

	t.Run("unmarshal_error", func(t *testing.T) {
		_, err := deploymentUnmarshal([]byte("a"))
		assert.Error(err)
	})
}

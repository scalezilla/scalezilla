package cluster

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/Lord-Y/rafty"
	"github.com/stretchr/testify/assert"
)

func TestFSM(t *testing.T) {
	assert := assert.New(t)

	t.Run("new_snapshot", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()

		snapshot := newSnapshot(cluster.config.DataDir, 3)
		assert.Nil(snapshot.List())
	})

	t.Run("snapshot_rafty_marshal_binary_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		fsm.raftyMarshalBinaryFunc = func(entry *rafty.LogEntry, w io.Writer) error {
			return errors.New("marshal binary error")
		}

		entry := &rafty.LogEntry{
			Term: 1,
		}
		fsm.memoryStore.logs[0] = entry

		buffer := new(bytes.Buffer)
		assert.Error(fsm.Snapshot(buffer))
	})

	t.Run("snapshot_rafty_marshal_binary_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		fsm.raftyMarshalBinaryFunc = func(entry *rafty.LogEntry, w io.Writer) error {
			return nil
		}
		fsm.raftyMarshalBinaryWithChecksumFunc = func(buffer *bytes.Buffer, w io.Writer) error {
			return errors.New("marshal binary checksum error")
		}

		entry := &rafty.LogEntry{
			Term: 1,
		}
		fsm.memoryStore.logs[0] = entry

		buffer := new(bytes.Buffer)
		assert.Error(fsm.Snapshot(buffer))
	})

	t.Run("snapshot_writer_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		fsm.raftyMarshalBinaryFunc = func(entry *rafty.LogEntry, w io.Writer) error {
			return nil
		}
		fsm.raftyMarshalBinaryWithChecksumFunc = func(buffer *bytes.Buffer, w io.Writer) error {
			return nil
		}

		entry := &rafty.LogEntry{
			Term: 1,
		}
		fsm.memoryStore.logs[0] = entry

		writer := &failWriter{failOn: 1}
		assert.Error(fsm.Snapshot(writer))
	})

	t.Run("snapshot_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		entry := &rafty.LogEntry{
			Term: 1,
		}
		fsm.memoryStore.logs[0] = entry

		buffer := new(bytes.Buffer)
		assert.Nil(fsm.Snapshot(buffer))
	})

	t.Run("apply_command_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		entry := &rafty.LogEntry{
			Term: 1,
		}
		fsm.memoryStoreLogsApplyCommandFunc = func(log *rafty.LogEntry) ([]byte, error) {
			return nil, errors.New("apply command error")
		}

		_, err = fsm.ApplyCommand(entry)
		assert.Error(err)
	})

	t.Run("apply_command_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		entry := &rafty.LogEntry{
			Term: 1,
		}

		_, err = fsm.ApplyCommand(entry)
		assert.Nil(err)
	})

	t.Run("apply_command_nil", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		entry := &rafty.LogEntry{
			Term:    1,
			LogType: uint32(rafty.LogReplication),
		}

		_, err = fsm.ApplyCommand(entry)
		assert.Nil(err)
	})

	t.Run("restore_eof_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		data := bytes.NewReader([]byte(`a=b`))
		assert.Error(fsm.Restore(data))
	})

	t.Run("restore_io_readfull_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		buffer := new(bytes.Buffer)
		assert.Nil(binary.Write(buffer, binary.LittleEndian, uint32(10)))
		_, err = buffer.Write([]byte{1, 2, 3}) // make insufficient data
		assert.Nil(err)
		data := bytes.NewReader(buffer.Bytes())
		assert.Error(fsm.Restore(data))
	})

	t.Run("restore_unmarshal_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		entry := &rafty.LogEntry{
			Term: 1,
		}

		fsm.memoryStore.logs[0] = entry
		fsm.raftyUnmarshalBinaryWithChecksumFunc = func(data []byte) (*rafty.LogEntry, error) {
			return nil, errors.New("unmarshal bianry error")
		}

		buffer := new(bytes.Buffer)
		assert.Nil(fsm.Snapshot(buffer))
		assert.Error(fsm.Restore(buffer))
	})

	t.Run("restore_apply_func_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		entry := &rafty.LogEntry{
			Term: 1,
		}

		fsm.memoryStore.logs[0] = entry
		fsm.raftyUnmarshalBinaryWithChecksumFunc = func(data []byte) (*rafty.LogEntry, error) {
			return nil, nil
		}
		fsm.applyCommandFunc = func(log *rafty.LogEntry) ([]byte, error) {
			return nil, errors.New("apply func error")
		}

		buffer := new(bytes.Buffer)
		assert.Nil(fsm.Snapshot(buffer))
		assert.Error(fsm.Restore(buffer))
	})

	t.Run("restore_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)

		fsm := newFSM(store)
		entry := &rafty.LogEntry{
			Term: 1,
		}

		fsm.memoryStore.logs[0] = entry

		buffer := new(bytes.Buffer)
		assert.Nil(fsm.Snapshot(buffer))
		assert.Nil(fsm.Restore(buffer))
	})
}

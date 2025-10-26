package cluster

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"maps"
	"slices"

	"github.com/Lord-Y/rafty"
)

// newFSM will return new fsm state
func newFSM(logStore rafty.LogStore) *fsmState {
	fsm := &fsmState{
		logStore: logStore,
		memoryStore: memoryStore{
			logs: make(map[uint64]*rafty.LogEntry),
		},
	}

	fsm.raftyMarshalBinaryFunc = rafty.MarshalBinary
	fsm.raftyMarshalBinaryWithChecksumFunc = rafty.MarshalBinaryWithChecksum
	fsm.raftyUnmarshalBinaryWithChecksumFunc = rafty.UnmarshalBinaryWithChecksum
	fsm.applyCommandFunc = fsm.ApplyCommand
	fsm.memoryStoreLogsApplyCommandFunc = fsm.memoryStore.logsApplyCommand

	return fsm
}

// newSnapshot will return new snapshot config
func newSnapshot(datadir string, maxSnapshot int) rafty.SnapshotStore {
	return rafty.NewSnapshot(datadir, maxSnapshot)
}

// Snapshot allows us to take snapshots
func (f *fsmState) Snapshot(snapshotWriter io.Writer) error {
	f.memoryStore.mu.RLock()
	defer f.memoryStore.mu.RUnlock()

	keys := slices.Sorted(maps.Keys(f.memoryStore.logs))
	for key := range keys {
		var err error
		buffer, bufferChecksum := new(bytes.Buffer), new(bytes.Buffer)
		if err = f.raftyMarshalBinaryFunc(f.memoryStore.logs[keys[key]], buffer); err != nil {
			return err
		}
		if err = f.raftyMarshalBinaryWithChecksumFunc(buffer, bufferChecksum); err != nil {
			return err
		}
		// writting data to the file handler
		if _, err = snapshotWriter.Write(bufferChecksum.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// Restore allows us to restore a snapshot
func (f *fsmState) Restore(snapshotReader io.Reader) error {
	reader := bufio.NewReader(snapshotReader)

	for {
		var length uint32
		if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Read first 4 bytes to get entry size
		record := make([]byte, length)
		if _, err := io.ReadFull(reader, record); err != nil {
			return err
		}

		data, err := f.raftyUnmarshalBinaryWithChecksumFunc(record)
		if err != nil {
			return err
		}
		if _, err := f.applyCommandFunc(data); err != nil {
			return err
		}
	}

	return nil
}

// ApplyCommand allows us to apply a command to the state machine.
func (f *fsmState) ApplyCommand(log *rafty.LogEntry) ([]byte, error) {
	if rafty.LogKind(log.LogType) == rafty.LogNoop || rafty.LogKind(log.LogType) == rafty.LogConfiguration {
		return f.memoryStoreLogsApplyCommandFunc(log)
	}

	return nil, nil
}

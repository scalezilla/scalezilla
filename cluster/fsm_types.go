package cluster

import (
	"bytes"
	"io"
	"sync"

	"github.com/Lord-Y/rafty"
)

// memoryStore holds the requirements related to user data
type memoryStore struct {
	// mu holds locking mecanism
	mu sync.RWMutex

	// logs map holds a map of the log entries
	logs map[uint64]*rafty.LogEntry
}

// fsmState is a struct holding a set of configs required for fsm
type fsmState struct {
	// LogStore is the store holding the data
	logStore rafty.LogStore

	// memoryStore is only for user land management
	memoryStore memoryStore

	// raftyMarshalBinaryFunc is used as a dependency injection
	raftyMarshalBinaryFunc func(entry *rafty.LogEntry, w io.Writer) error

	// raftyMarshalBinaryWithChecksumFunc is used as a dependency injection
	raftyMarshalBinaryWithChecksumFunc func(buffer *bytes.Buffer, w io.Writer) error

	// raftyUnmarshalBinaryWithChecksumFunc is used as a dependency injection
	raftyUnmarshalBinaryWithChecksumFunc func(data []byte) (*rafty.LogEntry, error)

	// applyCommandFunc is used as a dependency injection
	applyCommandFunc func(log *rafty.LogEntry) ([]byte, error)

	// memoryStoreLogsApplyCommandFunc is used as a dependency injection
	memoryStoreLogsApplyCommandFunc func(log *rafty.LogEntry) ([]byte, error)
}

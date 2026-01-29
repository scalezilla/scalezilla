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

	// aclToken map holds a map of decoded acl Token store
	aclToken map[string]dataACLToken
}

// data holds informations about the decoded command and its log index.
// this index is used to know what entry to drop or override
type dataACLToken struct {
	// index is the index of the log entry
	index uint64

	// value is the value of the command
	value aclTokenCommand
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

// commandKind represent the command that will be applied to the state machine
// It can be only be Get, Set, Delete
type commandKind uint32

const (
	// aclTokenCommandGet command allows us to fetch data from acl token fsm store
	aclTokenCommandGet commandKind = iota

	// aclTokenCommandGetAll command allows us to fetch all data from acl token fsm store
	aclTokenCommandGetAll

	// aclTokenCommandSet command allows us to write data from acl token fsm store
	aclTokenCommandSet

	// aclTokenCommandDelete command allows us to delete data from acl token fsm store
	aclTokenCommandDelete

	// dummyTest is only used during unit testing
	dummyTest
)

// aclTokenCommand is the struct to use to interact with cluster data
type aclTokenCommand struct {
	// Kind represent the set of commands: get, set, del
	Kind commandKind

	AccessorID string

	Token string

	InitialToken bool
}

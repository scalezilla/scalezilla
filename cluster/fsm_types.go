package cluster

import (
	"bytes"
	"io"
	"sync"
	"time"

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

	// deployment map holds a map of the decoded deployments
	deployment map[string]deploymentState
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

	// deploymentCommandGet command allows us to fetch data from deployment fsm store
	deploymentCommandGet

	// deploymentCommandGetAll command allows us to fetch all data from deployment fsm store
	deploymentCommandGetAll

	// deploymentCommandSet command allows us to write data from deployment fsm store
	deploymentCommandSet

	// deploymentCommandDelete command allows us to delete data from deployment fsm store
	deploymentCommandDelete
)

// aclTokenCommand is the struct to use to interact with cluster data
type aclTokenCommand struct {
	// Kind represent the set of commands: get, set, del
	Kind commandKind

	AccessorID string

	Token string

	InitialToken bool
}

// deploymentContent holds requirements related to deployments
type deploymentContent struct {
	//
	// IsStable tell if the deployment has successfully deployed
	IsStable bool `json:"is_stable"`

	// RawContent holds deployment content
	RawContent string `json:"raw_content"`

	// Version is the deployment version
	Version uint64 `json:"version"`

	// CreatedAt is the creation date
	CreatedAt time.Time `json:"created_at"`

	// ReplicaSetID is the name of the id of the replicaset
	ReplicaSetID string `json:"replicaset_id"`
}

// deploymentState is the state related to a deployment
type deploymentState struct {
	// Kind represent the set of commands: get, set, del
	Kind commandKind `json:"kind"`

	// index is the index of the log entry.
	// must not be used during encoding/decoding
	index uint64

	// Name is the deployment name
	Name string `json:"name"`

	// NewRollingVersion is used when a new version is ongoing.
	// set to -1 when no new version
	NewRollingVersion int64 `json:"new_rolling_version"`

	// CurrentUsedVersion is used to know which version to use
	// when is started after beeing shutdown
	CurrentUsedVersion uint64 `json:"current_used_version"`

	// Content holds deployment requirements
	Content map[uint64]deploymentContent `json:"content"`

	// MustBeStarted is a flag who by default is set to true.
	// When set to false, all pods will be stopped
	MustBeStarted bool `json:"must_be_started"`
}

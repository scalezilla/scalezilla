package cluster

import (
	"time"

	"github.com/Lord-Y/rafty"
)

// raftyStore is an interface implements Rafty requirements.
// This will be useful during unit testing
type raftyStore interface {
	// Close implementation of rafty bolt store
	Close() error

	// Close implementation of rafty bolt store
	CompactLogs(index uint64) error

	// Close implementation of rafty bolt store
	DiscardLogs(minIndex uint64, maxIndex uint64) error

	// Close implementation of rafty bolt store
	FirstIndex() (uint64, error)

	// Close implementation of rafty bolt store
	Get(key []byte) ([]byte, error)

	// Close implementation of rafty bolt store
	GetLastConfiguration() (*rafty.LogEntry, error)

	// Close implementation of rafty bolt store
	GetLogByIndex(index uint64) (*rafty.LogEntry, error)

	// Close implementation of rafty bolt store
	GetLogsByRange(minIndex uint64, maxIndex uint64, maxAppendEntries uint64) (response rafty.GetLogsByRangeResponse)

	// Close implementation of rafty bolt store
	GetMetadata() ([]byte, error)

	// Close implementation of rafty bolt store
	GetUint64(key []byte) uint64

	// Close implementation of rafty bolt store
	LastIndex() (uint64, error)

	// Close implementation of rafty bolt store
	Set(key []byte, value []byte) error

	// Close implementation of rafty bolt store
	SetUint64(key []byte, value []byte) error

	// Close implementation of rafty bolt store
	StoreLog(log *rafty.LogEntry) error

	// Close implementation of rafty bolt store
	StoreLogs(logs []*rafty.LogEntry) error

	// Close implementation of rafty bolt store
	StoreMetadata(value []byte) error
}

// raftyServer is an interface implements Rafty requirements.
// This will be useful during unit testing
type raftyServer interface {
	Start() error
	Stop()
	IsBootstrapped() bool
	SubmitCommand(timeout time.Duration, logKind rafty.LogKind, command []byte) ([]byte, error)
	BootstrapCluster(timeout time.Duration) error
	IsLeader() bool
	Leader() (bool, string, string)
	Status() rafty.Status
}

package cluster

import "errors"

var (
	// ErrServerOrClientModeMustUndefined is returned when there is NO
	// server block or client block defined.
	// End user must use one or them or both.
	ErrServerOrClientBlockUndefined = errors.New("server block or client block must be defined")

	// ErrServerRaftBlockInvalid is returned when client block is nil and server.enabled is false
	ErrServerRaftBlockInvalid = errors.New("server raft block invalid")

	// ErrClientRaftBlockInvalid is returned when server block is nil and client.enabled is false
	ErrClientRaftBlockInvalid = errors.New("client raft block invalid")

	ErrRaftBootstrapExpectedSizeInvalid = errors.New("raft bootstrap expected size cannot be zero")

	// ErrServerClusterJoinBlockInvalid is returned when cluster_join block is nil
	ErrServerClusterJoinBlockInvalid = errors.New("server cluster_join block invalid")

	// ErrClientClusterJoinBlockInvalid is returned when cluster_join block is nil
	ErrClientClusterJoinBlockInvalid = errors.New("client cluster_join block invalid")

	// ErrClusterJoinInitialMembersEmpty is returned when cluster_join.initial_members is empty
	ErrClusterJoinInitialMembersEmpty = errors.New("cluster_join.initial_members cannot be empty")

	// ErrClusterJoinInitialMembersInvalid is returned when parsed members
	// are invalid
	ErrClusterJoinInitialMembersInvalid = errors.New("parsed cluster_join.initial_members are invalid")

	// ErrServerOrClientModeMustBeEnabled is returned when server.enabled and client.enabled are both false.
	// User must set one of them to true
	ErrServerOrClientModeMustBeEnabled = errors.New("server or client flag 'enabled' must be set to true")

	// ErrServerSchedulerConfigBinpackMode is returned when server.scheduler_config.binpack_mode is invalid
	ErrServerSchedulerConfigBinpackMode = errors.New("server scheduler binpack mode can only be 'compact' or 'spread'")

	// ErrCgroupV2Required is returned when cgroup v2 is not detected on linux OS
	ErrCgroupV2Required = errors.New("cgroup v2 is required")

	// ErrOsUnsupported is returned when runtime.GOOS != "linux"
	// ErrOsUnsupported = errors.New("os unsupported")

	// ErrTimeout is triggered when an operation timed out
	ErrTimeout = errors.New("operation timeout")

	// ErrShutdown is triggered when the node is shutting down
	ErrShutdown = errors.New("node is shutting down")

	// ErrClusterNotBootstrapped is triggered when the cluster is not bootstrapped
	ErrClusterNotBootstrapped = errors.New("cluster not boostrapped")

	// ErrWrongFormat is triggered when the wrong format output is passed
	ErrWrongFormat = errors.New("wrong format")

	// ErrDeploymentPayload is triggered when the posting empty payload to deployments endpoints
	ErrDeploymentPayload = errors.New("deployment payload cannot be empty")

	// ErrDeploymentNameInvalid is triggered when deployment name is invalid
	ErrDeploymentNameInvalid = errors.New("deployment name must match ^[A-Za-z][A-Za-z0-9-]{5,62}")

	// ErrPodNameInvalid is triggered when pod name is invalid
	ErrPodNameInvalid = errors.New("pod name must match ^[A-Za-z][A-Za-z0-9-]{5,62}")

	// ErrContainerNameInvalid is triggered when container name is invalid
	ErrContainerNameInvalid = errors.New("container name must match ^[A-Za-z][A-Za-z0-9-]{5,62}")

	ErrPodsDeleteInvalid = errors.New("you must specify at least one container id")
)

// respError is used to decode errors
// from HTTP request
type respError struct {
	Error string `json:"error"`
}

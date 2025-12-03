package cluster

import "errors"

var (
	// ErrServerOrClientModeMustUndefined is return when there is NO
	// server block or client block defined.
	// End user must use one or them or both.
	ErrServerOrClientBlockUndefined = errors.New("server block or client block must be defined")

	// ErrServerRaftBlockInvalid is return when client block is nil and server.enabled is false
	ErrServerRaftBlockInvalid = errors.New("server raft block invalid")

	// ErrClientRaftBlockInvalid is return when server block is nil and client.enabled is false
	ErrClientRaftBlockInvalid = errors.New("client raft block invalid")

	ErrRaftBootstrapExpectedSizeInvalid = errors.New("raft bootstrap expected size cannot be zero")

	// ErrServerClusterJoinBlockInvalid is return when cluster_join block is nil
	ErrServerClusterJoinBlockInvalid = errors.New("server cluster_join block invalid")

	// ErrClientClusterJoinBlockInvalid is return when cluster_join block is nil
	ErrClientClusterJoinBlockInvalid = errors.New("client cluster_join block invalid")

	// ErrClusterJoinInitialMembersEmpty is returned when cluster_join.initial_members is empty
	ErrClusterJoinInitialMembersEmpty = errors.New("cluster_join.initial_members cannot be empty")

	// ErrClusterJoinInitialMembersInvalid is return when parsed members
	// are invalid
	ErrClusterJoinInitialMembersInvalid = errors.New("parsed cluster_join.initial_members are invalid")

	// ErrServerOrClientModeMustBeEnabled is return when server.enabled and client.enabled are both false.
	// User must set one of them to true
	ErrServerOrClientModeMustBeEnabled = errors.New("server or client flag 'enabled' must be set to true")

	// ErrCgroupV2Required is return when cgroup v2 is not detected on linux OS
	ErrCgroupV2Required = errors.New("cgroup v2 is required")

	// ErrOsUnsupported is returned when runtime.GOOS != "linux"
	// ErrOsUnsupported = errors.New("os unsupported")

	// ErrTimeout is triggered when an operation timed out
	ErrTimeout = errors.New("operation timeout")

	// ErrShutdown is triggered when the node is shutting down
	ErrShutdown = errors.New("node is shutting down")
)

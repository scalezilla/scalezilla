package cri

import (
	"context"

	tasksapi "github.com/containerd/containerd/api/services/tasks/v1"
	containerd "github.com/containerd/containerd/v2/client"
)

// containerdImage wraps a pulled containerd image for the runtime seam.
type containerdImage struct {
	image containerd.Image
}

// containerdClientAPI captures the containerd client methods used by the adapter.
type containerdClientAPI interface {
	Pull(ctx context.Context, ref string, opts ...containerd.RemoteOpt) (containerd.Image, error)
	NewContainer(ctx context.Context, id string, opts ...containerd.NewContainerOpts) (containerd.Container, error)
	Containers(ctx context.Context, filters ...string) ([]containerd.Container, error)
	TaskService() tasksapi.TasksClient
	LoadContainer(ctx context.Context, id string) (containerd.Container, error)
	Close() error
}

// containerdRuntimeClient adapts the containerd client to runtimeClient.
type containerdRuntimeClient struct {
	client containerdClientAPI
}

// containerdClientFactory creates the concrete containerd client for the adapter.
var containerdClientFactory = func(address string) (containerdClientAPI, error) {
	return containerd.New(address)
}

// containerdRuntimeContainer adapts containerd.Container to runtimeContainer.
type containerdRuntimeContainer struct {
	container containerd.Container
}

// containerdRuntimeTask adapts containerd.Task to runtimeTask.
type containerdRuntimeTask struct {
	task containerd.Task
}

// containerdRuntimeExitStatus wraps containerd.ExitStatus for runtimeExitStatus.
type containerdRuntimeExitStatus struct {
	status containerd.ExitStatus
}

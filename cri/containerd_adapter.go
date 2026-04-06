package cri

import (
	"context"
	"fmt"
	"syscall"
	"time"

	tasksapi "github.com/containerd/containerd/api/services/tasks/v1"
	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/cio"
	"github.com/containerd/containerd/v2/pkg/oci"
)

// isRuntimeImage marks containerdImage as satisfying runtimeImage.
func (containerdImage) isRuntimeImage() {}

// newContainerdClient builds the default runtimeClient backed by containerd.
func newContainerdClient() (runtimeClient, error) {
	client, err := containerdClientFactory(containerdAddress)
	if err != nil {
		return nil, err
	}
	return &containerdRuntimeClient{client: client}, nil
}

// Pull retrieves and unpacks an image using the underlying containerd client.
func (c *containerdRuntimeClient) Pull(ctx context.Context, ref string) (runtimeImage, error) {
	image, err := c.client.Pull(ctx, ref, containerd.WithPullUnpack)
	if err != nil {
		return nil, err
	}
	return containerdImage{image: image}, nil
}

// NewContainer creates a containerd container and wraps it for runtimeClient use.
func (c *containerdRuntimeClient) NewContainer(ctx context.Context, id string, image runtimeImage, labels, additionalContainerLabels map[string]string) (runtimeContainer, error) {
	containerdImage, ok := image.(containerdImage)
	if !ok {
		return nil, fmt.Errorf("unexpected runtime image type %T", image)
	}

	container, err := c.client.NewContainer(
		ctx,
		id,
		containerd.WithImage(containerdImage.image),
		containerd.WithNewSnapshot(fmt.Sprintf("snapshot-%s", id), containerdImage.image),
		containerd.WithNewSpec(oci.WithImageConfig(containerdImage.image)),
		containerd.WithContainerLabels(labels),
		containerd.WithAdditionalContainerLabels(additionalContainerLabels),
	)
	if err != nil {
		return nil, err
	}

	return &containerdRuntimeContainer{container: container}, nil
}

// Containers lists containers and converts them to runtimeContainer values.
func (c *containerdRuntimeClient) Containers(ctx context.Context) ([]runtimeContainer, error) {
	containers, err := c.client.Containers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]runtimeContainer, 0, len(containers))
	for _, container := range containers {
		result = append(result, &containerdRuntimeContainer{container: container})
	}
	return result, nil
}

// ListTasks lists containerd tasks and maps them to the reduced runtime shape.
func (c *containerdRuntimeClient) ListTasks(ctx context.Context) ([]runtimeTaskProcess, error) {
	taskResp, err := c.client.TaskService().List(ctx, &tasksapi.ListTasksRequest{})
	if err != nil {
		return nil, err
	}

	result := make([]runtimeTaskProcess, 0, len(taskResp.Tasks))
	for _, task := range taskResp.Tasks {
		result = append(result, runtimeTaskProcess{
			ID:     task.ID,
			PID:    task.Pid,
			Status: task.Status.String(),
		})
	}
	return result, nil
}

// LoadContainer loads a container by ID and wraps it for runtimeContainer use.
func (c *containerdRuntimeClient) LoadContainer(ctx context.Context, id string) (runtimeContainer, error) {
	container, err := c.client.LoadContainer(ctx, id)
	if err != nil {
		return nil, err
	}
	return &containerdRuntimeContainer{container: container}, nil
}

// Close closes the underlying containerd client connection.
func (c *containerdRuntimeClient) Close() error {
	return c.client.Close()
}

// ID returns the wrapped container identifier.
func (c *containerdRuntimeContainer) ID() string {
	return c.container.ID()
}

// Info returns the container fields needed by ListContainer.
func (c *containerdRuntimeContainer) Info(ctx context.Context) (runtimeContainerInfo, error) {
	info, err := c.container.Info(ctx, containerd.WithoutRefreshedMetadata)
	if err != nil {
		return runtimeContainerInfo{}, err
	}

	return runtimeContainerInfo{
		Image:     info.Image,
		Runtime:   info.Runtime.Name,
		Labels:    info.Labels,
		CreatedAt: info.CreatedAt,
	}, nil
}

// NewTask creates a task configured to write logs to the provided file.
func (c *containerdRuntimeContainer) NewTask(ctx context.Context, logFile string) (runtimeTask, error) {
	task, err := c.container.NewTask(ctx, cio.LogFile(logFile))
	if err != nil {
		return nil, err
	}
	return &containerdRuntimeTask{task: task}, nil
}

// Task loads the current task for the wrapped container.
func (c *containerdRuntimeContainer) Task(ctx context.Context) (runtimeTask, error) {
	task, err := c.container.Task(ctx, cio.Load)
	if err != nil {
		return nil, err
	}
	return &containerdRuntimeTask{task: task}, nil
}

// Delete removes the wrapped container and cleans up its snapshot.
func (c *containerdRuntimeContainer) Delete(ctx context.Context) error {
	return c.container.Delete(ctx, containerd.WithSnapshotCleanup)
}

// StopSignal resolves the stop signal configured for the wrapped container.
func (c *containerdRuntimeContainer) StopSignal(ctx context.Context, defaultSignal syscall.Signal) (syscall.Signal, error) {
	return containerd.GetStopSignal(ctx, c.container, defaultSignal)
}

// ID returns the wrapped task identifier.
func (t *containerdRuntimeTask) ID() string {
	return t.task.ID()
}

// Wait bridges containerd exit statuses into the runtimeExitStatus interface.
func (t *containerdRuntimeTask) Wait(ctx context.Context) (<-chan runtimeExitStatus, error) {
	waitC, err := t.task.Wait(ctx)
	if err != nil {
		return nil, err
	}

	result := make(chan runtimeExitStatus, 1)
	go func() {
		defer close(result)
		exitStatus, ok := <-waitC
		if !ok {
			result <- containerdRuntimeExitStatus{}
			return
		}
		result <- containerdRuntimeExitStatus{status: exitStatus}
	}()

	return result, nil
}

// Start starts the wrapped task.
func (t *containerdRuntimeTask) Start(ctx context.Context) error {
	return t.task.Start(ctx)
}

// Status reads the wrapped task status and converts it to runtimeTaskStatus.
func (t *containerdRuntimeTask) Status(ctx context.Context) (runtimeTaskStatus, error) {
	status, err := t.task.Status(ctx)
	if err != nil {
		return "", err
	}
	return runtimeTaskStatus(status.Status), nil
}

// Kill sends the provided signal to the wrapped task.
func (t *containerdRuntimeTask) Kill(ctx context.Context, sig syscall.Signal) error {
	return t.task.Kill(ctx, sig)
}

// Delete removes the wrapped task and returns its exit status.
func (t *containerdRuntimeTask) Delete(ctx context.Context) (runtimeExitStatus, error) {
	exitStatus, err := t.task.Delete(ctx, containerd.WithProcessKill)
	if err != nil {
		return nil, err
	}
	if exitStatus == nil {
		return containerdRuntimeExitStatus{}, nil
	}
	return containerdRuntimeExitStatus{status: *exitStatus}, nil
}

// Result returns the code, exit time, and any wait error from the wrapped status.
func (s containerdRuntimeExitStatus) Result() (uint32, time.Time, error) {
	return s.status.Result()
}

package cri

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/namespaces"
	"github.com/containerd/errdefs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewCRI returns an interface to manage containers
func NewCRI(log *zerolog.Logger) ContainerRuntime {
	return &CRI{
		log:       log,
		newClient: newContainerdClient,
		mkdirAll:  os.MkdirAll,
		after:     time.After,
	}
}

// client returns the injected runtime client or the default containerd-backed one.
func (c *CRI) client() (runtimeClient, error) {
	if c.newClient != nil {
		return c.newClient()
	}
	return newContainerdClient()
}

// makeDir delegates directory creation to the injected filesystem hook when present.
func (c *CRI) makeDir(path string, perm os.FileMode) error {
	if c.mkdirAll != nil {
		return c.mkdirAll(path, perm)
	}
	return os.MkdirAll(path, perm)
}

// afterSignal returns the injected timer channel or falls back to time.After.
func (c *CRI) afterSignal(timeout time.Duration) <-chan time.Time {
	if c.after != nil {
		return c.after(timeout)
	}
	return time.After(timeout)
}

// CreateContainer creates container with its requirements
func (c *CRI) CreateContainer(ctx context.Context, spec CreateContainerSpec) error {
	client, err := c.client()
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", spec.Namespace).
			Str("container", spec.ContainerID).
			Msg("Fail to connect to containerd socket")
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	cctx := namespaces.WithNamespace(ctx, spec.Namespace)
	image, err := client.Pull(cctx, spec.Image.Image)
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", spec.Namespace).
			Str("container", spec.ContainerID).
			Msg("Fail to pull image")
		return err
	}

	container, err := client.NewContainer(cctx, spec.ContainerID, image)
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", spec.Namespace).
			Str("container", spec.ContainerID).
			Msg("Fail to create container config")
		return err
	}

	containerLogFile := filepath.Join(spec.DefaultLogPath, spec.Namespace, spec.ContainerID, fmt.Sprintf("%s.log", spec.ContainerID))
	if err := c.makeDir(filepath.Dir(containerLogFile), 0755); err != nil {
		c.log.Error().Err(err).
			Str("namespace", spec.Namespace).
			Str("container", spec.ContainerID).
			Msg("Fail to create log directory")
		return err
	}

	task, err := container.NewTask(cctx, containerLogFile)
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", spec.Namespace).
			Str("container", spec.ContainerID).
			Msg("Fail to create container task")
		return err
	}

	if _, err = task.Wait(cctx); err != nil {
		c.log.Error().Err(err).
			Str("namespace", spec.Namespace).
			Str("container", spec.ContainerID).
			Msg("Fail to wait container task")
		return err
	}

	if err := task.Start(cctx); err != nil {
		c.log.Error().Err(err).
			Str("namespace", spec.Namespace).
			Str("container", spec.ContainerID).
			Msg("Fail to start container task")
		return err
	}

	log.Debug().
		Str("namespace", spec.Namespace).
		Str("container", spec.ContainerID).
		Msg("Container started")

	return nil
}

// ListContainer returns the container list
func (c *CRI) ListContainer(ctx context.Context, namespace string) ([]ContainerList, error) {
	client, err := c.client()
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", namespace).
			Msg("Fail to connect to containerd socket")
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()

	cctx := namespaces.WithNamespace(ctx, namespace)
	containers, err := client.Containers(cctx)
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", namespace).
			Msg("Fail to list container")
		return nil, err
	}

	taskList, err := client.ListTasks(cctx)
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", namespace).
			Msg("Fail to list tasks")
		return nil, err
	}

	taskByID := make(map[string]runtimeTaskProcess, len(taskList))
	for _, task := range taskList {
		taskByID[task.ID] = task
	}

	var list []ContainerList
	for _, cc := range containers {
		info, err := cc.Info(cctx)
		if err != nil {
			c.log.Error().Err(err).
				Str("namespace", namespace).
				Str("container", cc.ID()).
				Msg("Fail to get container info")
		}

		image := info.Image
		if image == "" {
			image = "-"
		}

		runtimeName := info.Runtime
		if runtimeName == "" {
			runtimeName = "-"
		}

		pid := uint32(0)
		status := "-"
		if task, ok := taskByID[cc.ID()]; ok {
			pid = task.PID
			status = task.Status
		}

		list = append(list, ContainerList{
			Namespace: namespace,
			ID:        cc.ID(),
			PID:       pid,
			Image:     image,
			Runtime:   runtimeName,
			Status:    status,
		})
	}

	return list, nil
}

// DeleteContainer deletes the provided container
func (c *CRI) DeleteContainer(ctx context.Context, namespace, containerID string, stopTimeout time.Duration) error {
	client, err := c.client()
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", namespace).
			Str("container", containerID).
			Msg("Fail to connect to containerd socket")
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	cctx := namespaces.WithNamespace(ctx, namespace)
	container, err := client.LoadContainer(cctx, containerID)
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", namespace).
			Str("container", containerID).
			Msg("Fail to load container metadata")
		return err
	}

	task, err := container.Task(cctx)
	switch {
	case err == nil:
		if err := c.stopAndDeleteTask(cctx, namespace, container, task, stopTimeout); err != nil {
			return err
		}
	case errors.Is(err, errdefs.ErrNotFound):
		// no task exists anymore; continue with container deletion
	default:
		c.log.Error().Err(err).
			Str("namespace", namespace).
			Str("container", containerID).
			Msg("Fail to load container task")
		return err
	}

	if err := container.Delete(cctx); err != nil {
		if errors.Is(err, errdefs.ErrNotFound) {
			c.log.Debug().
				Str("namespace", namespace).
				Str("container", containerID).
				Msg("Container already deleted")
			return nil
		}

		c.log.Error().Err(err).
			Str("namespace", namespace).
			Str("container", containerID).
			Msg("Fail to delete container")
		return err
	}

	return nil
}

// StopAndDeleteTask is stopping and deleting the task related to the provided container
func (c *CRI) StopAndDeleteTask(ctx context.Context, namespace string, container containerd.Container, task containerd.Task, stopTimeout time.Duration) error {
	return c.stopAndDeleteTask(ctx, namespace, &containerdRuntimeContainer{container: container}, &containerdRuntimeTask{task: task}, stopTimeout)
}

// stopAndDeleteTask contains the shutdown flow shared by real and fake runtime types.
func (c *CRI) stopAndDeleteTask(ctx context.Context, namespace string, container runtimeContainer, task runtimeTask, stopTimeout time.Duration) error {
	status, err := task.Status(ctx)
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", namespace).
			Str("container", container.ID()).
			Str("task", task.ID()).
			Msg("Fail to get task status")
		return err
	}

	if status != runtimeTaskStopped {
		waitC, err := task.Wait(ctx)
		if err != nil {
			c.log.Error().Err(err).
				Str("namespace", namespace).
				Str("container", container.ID()).
				Str("task", task.ID()).
				Msg("Fail to wait for task")
			return err
		}

		sig, err := container.StopSignal(ctx, syscall.SIGTERM)
		if err != nil {
			sig = syscall.SIGTERM
		}

		if err := task.Kill(ctx, sig); err != nil && !errors.Is(err, errdefs.ErrNotFound) {
			c.log.Error().Err(err).
				Str("namespace", namespace).
				Str("container", container.ID()).
				Str("task", task.ID()).
				Msgf("Fail to kill signal with %v", sig)
			return err
		}

		select {
		case exitStatus := <-waitC:
			code, exitedAt, waitErr := exitStatus.Result()
			if waitErr != nil {
				c.log.Error().Err(waitErr).
					Str("namespace", namespace).
					Str("container", container.ID()).
					Str("task", task.ID()).
					Msg("Read exit status is wrong")
				return waitErr
			}

			c.log.Debug().
				Str("namespace", namespace).
				Str("container", container.ID()).
				Str("task", task.ID()).
				Str("code", fmt.Sprintf("%d", code)).
				Str("at", exitedAt.Format(time.RFC3339)).
				Msg("Task exited gracefully")

		case <-c.afterSignal(stopTimeout):
			c.log.Info().
				Str("namespace", namespace).
				Str("container", container.ID()).
				Str("task", task.ID()).
				Msgf("Task did not stop within %s, sending SIGKILL", stopTimeout)

			if err := task.Kill(ctx, syscall.SIGKILL); err != nil && !errors.Is(err, errdefs.ErrNotFound) {
				c.log.Error().Err(err).
					Str("namespace", namespace).
					Str("container", container.ID()).
					Str("task", task.ID()).
					Msg("Fail to SIGKILL task")
				return err
			}

			exitStatus := <-waitC
			code, exitedAt, waitErr := exitStatus.Result()
			if waitErr != nil {
				c.log.Error().Err(waitErr).
					Str("namespace", namespace).
					Str("container", container.ID()).
					Str("task", task.ID()).
					Msg("Fail to read exit status after SIGKILL")
				return waitErr
			}

			c.log.Debug().
				Str("namespace", namespace).
				Str("container", container.ID()).
				Str("task", task.ID()).
				Str("code", fmt.Sprintf("%d", code)).
				Str("at", exitedAt.Format(time.RFC3339)).
				Msg("Task exited after SIGKILL")
		}
	}

	exitStatus, err := task.Delete(ctx)
	if err != nil {
		if errors.Is(err, errdefs.ErrNotFound) {
			c.log.Debug().
				Str("namespace", namespace).
				Str("container", container.ID()).
				Str("task", task.ID()).
				Msg("Task already deleted")
			return nil
		}

		c.log.Error().Err(err).
			Str("namespace", namespace).
			Str("container", container.ID()).
			Str("task", task.ID()).
			Msg("Fail to delete task")
		return err
	}

	code, exitedAt, err := exitStatus.Result()
	if err != nil {
		c.log.Error().Err(err).
			Str("namespace", namespace).
			Str("container", container.ID()).
			Str("task", task.ID()).
			Msg("Fail to delete read deleted task exit status")
	}

	c.log.Debug().
		Str("namespace", namespace).
		Str("container", container.ID()).
		Str("task", task.ID()).
		Str("code", fmt.Sprintf("%d", code)).
		Str("at", exitedAt.Format(time.RFC3339)).
		Msg("Task deleted")
	return nil
}

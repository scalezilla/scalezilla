package cri

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/containerd/errdefs"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewCRI(t *testing.T) {
	_, ok := NewCRI(testLogger()).(*CRI)
	require.True(t, ok)
}

func TestCRI_CreateContainer(t *testing.T) {
	spec := CreateContainerSpec{
		Namespace:   "example",
		ContainerID: "redis-server",
		Image: ImageSpec{
			Image: "docker.io/library/redis:alpine",
		},
		DefaultLogPath: "/tmp/containerd",
	}

	t.Run("client error", func(t *testing.T) {
		wantErr := errors.New("dial failed")
		cri := &CRI{
			log: testLogger(),
			newClient: func() (runtimeClient, error) {
				return nil, wantErr
			},
		}

		err := cri.CreateContainer(context.Background(), spec)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("pull error", func(t *testing.T) {
		wantErr := errors.New("pull failed")
		client := &fakeClient{pullErr: wantErr}
		cri := newTestCRI(client)

		err := cri.CreateContainer(context.Background(), spec)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("new container error", func(t *testing.T) {
		wantErr := errors.New("new container failed")
		client := &fakeClient{
			pullImage:       fakeImage{},
			newContainerErr: wantErr,
		}
		cri := newTestCRI(client)

		err := cri.CreateContainer(context.Background(), spec)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("mkdir error", func(t *testing.T) {
		wantErr := errors.New("mkdir failed")
		task := &fakeTask{id: "redis-server", waitCh: make(chan runtimeExitStatus, 1)}
		container := &fakeContainer{id: spec.ContainerID, newTask: task}
		client := &fakeClient{
			pullImage:    fakeImage{},
			newContainer: container,
		}

		cri := newTestCRI(client)
		cri.mkdirAll = func(path string, perm os.FileMode) error {
			return wantErr
		}

		err := cri.CreateContainer(context.Background(), spec)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("new task error", func(t *testing.T) {
		wantErr := errors.New("new task failed")
		container := &fakeContainer{id: spec.ContainerID, newTaskErr: wantErr}
		client := &fakeClient{
			pullImage:    fakeImage{},
			newContainer: container,
		}

		cri := newTestCRI(client)

		err := cri.CreateContainer(context.Background(), spec)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("wait error", func(t *testing.T) {
		wantErr := errors.New("wait failed")
		task := &fakeTask{id: "redis-server", waitErr: wantErr}
		container := &fakeContainer{id: spec.ContainerID, newTask: task}
		client := &fakeClient{
			pullImage:    fakeImage{},
			newContainer: container,
		}

		cri := newTestCRI(client)

		err := cri.CreateContainer(context.Background(), spec)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("start error", func(t *testing.T) {
		wantErr := errors.New("start failed")
		task := &fakeTask{
			id:       "redis-server",
			waitCh:   immediateWait(fakeExitStatus{}),
			startErr: wantErr,
		}
		container := &fakeContainer{id: spec.ContainerID, newTask: task}
		client := &fakeClient{
			pullImage:    fakeImage{},
			newContainer: container,
		}

		cri := newTestCRI(client)

		err := cri.CreateContainer(context.Background(), spec)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("success", func(t *testing.T) {
		task := &fakeTask{
			id:     "redis-server",
			waitCh: immediateWait(fakeExitStatus{}),
		}
		container := &fakeContainer{id: spec.ContainerID, newTask: task}
		client := &fakeClient{
			pullImage:    fakeImage{},
			newContainer: container,
		}

		var createdDir string
		cri := newTestCRI(client)
		cri.mkdirAll = func(path string, perm os.FileMode) error {
			createdDir = path
			return nil
		}

		err := cri.CreateContainer(context.Background(), spec)
		require.NoError(t, err)
		require.True(t, client.closed)
		require.Equal(t, spec.Image.Image, client.pullRef)
		require.Equal(t, spec.ContainerID, client.newContainerID)
		require.Equal(t, filepath.Join(spec.DefaultLogPath, spec.Namespace, spec.ContainerID), createdDir)
		require.Equal(t, filepath.Join(spec.DefaultLogPath, spec.Namespace, spec.ContainerID, "redis-server.log"), container.newTaskLogFile)
		require.True(t, task.started)
	})
}

func TestCRI_ListContainer(t *testing.T) {
	t.Run("client error", func(t *testing.T) {
		wantErr := errors.New("dial failed")
		cri := &CRI{
			log: testLogger(),
			newClient: func() (runtimeClient, error) {
				return nil, wantErr
			},
		}

		_, err := cri.ListContainer(context.Background(), "example")
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("containers error", func(t *testing.T) {
		wantErr := errors.New("list failed")
		client := &fakeClient{containersErr: wantErr}
		cri := newTestCRI(client)

		_, err := cri.ListContainer(context.Background(), "example")
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("tasks error", func(t *testing.T) {
		wantErr := errors.New("tasks failed")
		client := &fakeClient{
			containers:   []runtimeContainer{&fakeContainer{id: "c1"}},
			listTasksErr: wantErr,
		}
		cri := newTestCRI(client)

		_, err := cri.ListContainer(context.Background(), "example")
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("maps container info and defaults", func(t *testing.T) {
		client := &fakeClient{
			containers: []runtimeContainer{
				&fakeContainer{
					id:   "c1",
					info: runtimeContainerInfo{Image: "redis:alpine", Runtime: "runc"},
				},
				&fakeContainer{
					id:      "c2",
					infoErr: errors.New("info failed"),
				},
			},
			taskList: []runtimeTaskProcess{
				{ID: "c1", PID: 42, Status: "running"},
			},
		}
		cri := newTestCRI(client)

		list, err := cri.ListContainer(context.Background(), "example")
		require.NoError(t, err)
		require.True(t, client.closed)
		require.Equal(t, []ContainerList{
			{
				Namespace: "example",
				ID:        "c1",
				PID:       42,
				Image:     "redis:alpine",
				Runtime:   "runc",
				Status:    "running",
			},
			{
				Namespace: "example",
				ID:        "c2",
				PID:       0,
				Image:     "-",
				Runtime:   "-",
				Status:    "-",
			},
		}, list)
	})
}

func TestCRI_DeleteContainer(t *testing.T) {
	t.Run("client error", func(t *testing.T) {
		wantErr := errors.New("dial failed")
		cri := &CRI{
			log: testLogger(),
			newClient: func() (runtimeClient, error) {
				return nil, wantErr
			},
		}

		err := cri.DeleteContainer(context.Background(), "example", "c1", time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("load container error", func(t *testing.T) {
		wantErr := errors.New("load failed")
		client := &fakeClient{loadContainerErr: wantErr}
		cri := newTestCRI(client)

		err := cri.DeleteContainer(context.Background(), "example", "c1", time.Second)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("task lookup error", func(t *testing.T) {
		wantErr := errors.New("task failed")
		container := &fakeContainer{id: "c1", taskErr: wantErr}
		client := &fakeClient{loadContainer: container}
		cri := newTestCRI(client)

		err := cri.DeleteContainer(context.Background(), "example", "c1", time.Second)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("task not found still deletes container", func(t *testing.T) {
		container := &fakeContainer{id: "c1", taskErr: errdefs.ErrNotFound}
		client := &fakeClient{loadContainer: container}
		cri := newTestCRI(client)

		err := cri.DeleteContainer(context.Background(), "example", "c1", time.Second)
		require.NoError(t, err)
		require.True(t, client.closed)
		require.True(t, container.deleted)
	})

	t.Run("stop and delete task error", func(t *testing.T) {
		wantErr := errors.New("status failed")
		task := &fakeTask{id: "c1", statusErr: wantErr}
		container := &fakeContainer{id: "c1", task: task}
		client := &fakeClient{loadContainer: container}
		cri := newTestCRI(client)

		err := cri.DeleteContainer(context.Background(), "example", "c1", time.Second)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
		require.False(t, container.deleted)
	})

	t.Run("container already deleted", func(t *testing.T) {
		task := &fakeTask{id: "c1", status: runtimeTaskStopped, deleteExitStatus: fakeExitStatus{}}
		container := &fakeContainer{id: "c1", task: task, deleteErr: errdefs.ErrNotFound}
		client := &fakeClient{loadContainer: container}
		cri := newTestCRI(client)

		err := cri.DeleteContainer(context.Background(), "example", "c1", time.Second)
		require.NoError(t, err)
		require.True(t, client.closed)
	})

	t.Run("container delete error", func(t *testing.T) {
		wantErr := errors.New("delete failed")
		task := &fakeTask{id: "c1", status: runtimeTaskStopped, deleteExitStatus: fakeExitStatus{}}
		container := &fakeContainer{id: "c1", task: task, deleteErr: wantErr}
		client := &fakeClient{loadContainer: container}
		cri := newTestCRI(client)

		err := cri.DeleteContainer(context.Background(), "example", "c1", time.Second)
		require.ErrorIs(t, err, wantErr)
		require.True(t, client.closed)
	})

	t.Run("success", func(t *testing.T) {
		task := &fakeTask{id: "c1", status: runtimeTaskStopped, deleteExitStatus: fakeExitStatus{}}
		container := &fakeContainer{id: "c1", task: task}
		client := &fakeClient{loadContainer: container}
		cri := newTestCRI(client)

		err := cri.DeleteContainer(context.Background(), "example", "c1", time.Second)
		require.NoError(t, err)
		require.True(t, client.closed)
		require.True(t, task.deleted)
		require.True(t, container.deleted)
	})
}

func TestCRI_StopAndDeleteTask(t *testing.T) {
	t.Run("status error", func(t *testing.T) {
		wantErr := errors.New("status failed")
		cri := newTestCRI(nil)

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1"}, &fakeTask{id: "t1", statusErr: wantErr}, time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("wait error", func(t *testing.T) {
		wantErr := errors.New("wait failed")
		cri := newTestCRI(nil)

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1"}, &fakeTask{
			id:      "t1",
			status:  runtimeTaskStatus("running"),
			waitErr: wantErr,
		}, time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("kill error", func(t *testing.T) {
		wantErr := errors.New("kill failed")
		cri := newTestCRI(nil)

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1", stopSignal: syscall.SIGUSR1}, &fakeTask{
			id:       "t1",
			status:   runtimeTaskStatus("running"),
			waitCh:   make(chan runtimeExitStatus, 1),
			killErrs: []error{wantErr},
		}, time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("graceful stop uses stop signal and tolerates not found kill", func(t *testing.T) {
		waitCh := make(chan runtimeExitStatus, 1)
		waitCh <- fakeExitStatus{code: 0, at: time.Unix(10, 0)}
		task := &fakeTask{
			id:               "t1",
			status:           runtimeTaskStatus("running"),
			waitCh:           waitCh,
			killErrs:         []error{errdefs.ErrNotFound},
			deleteExitStatus: fakeExitStatus{},
		}
		container := &fakeContainer{id: "c1", stopSignal: syscall.SIGUSR1}
		cri := newTestCRI(nil)

		err := cri.stopAndDeleteTask(context.Background(), "example", container, task, time.Second)
		require.NoError(t, err)
		require.Equal(t, []syscall.Signal{syscall.SIGUSR1}, task.killedSignals)
		require.True(t, task.deleted)
	})

	t.Run("graceful exit result error bubbles up", func(t *testing.T) {
		wantErr := errors.New("result failed")
		waitCh := make(chan runtimeExitStatus, 1)
		waitCh <- fakeExitStatus{err: wantErr}
		cri := newTestCRI(nil)
		container := &fakeContainer{id: "c1", stopSignalErr: errors.New("lookup failed")}

		err := cri.stopAndDeleteTask(context.Background(), "example", container, &fakeTask{
			id:     "t1",
			status: runtimeTaskStatus("running"),
			waitCh: waitCh,
		}, time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("timeout sends sigkill", func(t *testing.T) {
		waitCh := make(chan runtimeExitStatus, 1)
		task := &fakeTask{
			id:               "t1",
			status:           runtimeTaskStatus("running"),
			waitCh:           waitCh,
			deleteExitStatus: fakeExitStatus{},
			onKill: func(sig syscall.Signal) {
				if sig == syscall.SIGKILL {
					waitCh <- fakeExitStatus{code: 137, at: time.Unix(20, 0)}
				}
			},
		}
		cri := newTestCRI(nil)
		cri.after = immediateAfter()

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1"}, task, time.Second)
		require.NoError(t, err)
		require.Equal(t, []syscall.Signal{syscall.SIGTERM, syscall.SIGKILL}, task.killedSignals)
		require.True(t, task.deleted)
	})

	t.Run("sigkill error", func(t *testing.T) {
		wantErr := errors.New("sigkill failed")
		waitCh := make(chan runtimeExitStatus, 1)
		task := &fakeTask{
			id:       "t1",
			status:   runtimeTaskStatus("running"),
			waitCh:   waitCh,
			killErrs: []error{nil, wantErr},
		}
		cri := newTestCRI(nil)
		cri.after = immediateAfter()

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1"}, task, time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("timeout exit result error bubbles up", func(t *testing.T) {
		wantErr := errors.New("result failed")
		waitCh := make(chan runtimeExitStatus, 1)
		task := &fakeTask{
			id:     "t1",
			status: runtimeTaskStatus("running"),
			waitCh: waitCh,
			onKill: func(sig syscall.Signal) {
				if sig == syscall.SIGKILL {
					waitCh <- fakeExitStatus{err: wantErr}
				}
			},
		}
		cri := newTestCRI(nil)
		cri.after = immediateAfter()

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1"}, task, time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("task already deleted", func(t *testing.T) {
		task := &fakeTask{id: "t1", status: runtimeTaskStopped, deleteErr: errdefs.ErrNotFound}
		cri := newTestCRI(nil)

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1"}, task, time.Second)
		require.NoError(t, err)
	})

	t.Run("task delete error", func(t *testing.T) {
		wantErr := errors.New("delete failed")
		task := &fakeTask{id: "t1", status: runtimeTaskStopped, deleteErr: wantErr}
		cri := newTestCRI(nil)

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1"}, task, time.Second)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("delete exit result error is logged and ignored", func(t *testing.T) {
		task := &fakeTask{
			id:               "t1",
			status:           runtimeTaskStopped,
			deleteExitStatus: fakeExitStatus{err: errors.New("result failed")},
		}
		cri := newTestCRI(nil)

		err := cri.stopAndDeleteTask(context.Background(), "example", &fakeContainer{id: "c1"}, task, time.Second)
		require.NoError(t, err)
		require.True(t, task.deleted)
	})
}

func newTestCRI(client runtimeClient) *CRI {
	return &CRI{
		log: testLogger(),
		newClient: func() (runtimeClient, error) {
			return client, nil
		},
		mkdirAll: func(string, os.FileMode) error {
			return nil
		},
		after: func(time.Duration) <-chan time.Time {
			return make(chan time.Time)
		},
	}
}

func testLogger() *zerolog.Logger {
	logger := zerolog.New(io.Discard)
	return &logger
}

func immediateAfter() func(time.Duration) <-chan time.Time {
	return func(time.Duration) <-chan time.Time {
		ch := make(chan time.Time, 1)
		ch <- time.Time{}
		return ch
	}
}

func immediateWait(status runtimeExitStatus) chan runtimeExitStatus {
	ch := make(chan runtimeExitStatus, 1)
	ch <- status
	return ch
}

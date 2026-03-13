package cri

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	tasksapi "github.com/containerd/containerd/api/services/tasks/v1"
	tasktypes "github.com/containerd/containerd/api/types/task"
	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/core/containers"
	"github.com/stretchr/testify/require"
)

func TestCRI_HelperFallbacks(t *testing.T) {
	t.Run("client fallback uses containerd factory", func(t *testing.T) {
		restore := swapContainerdClientFactory(func(string) (containerdClientAPI, error) {
			return &fakeContainerdClientAPI{}, nil
		})
		defer restore()

		cri := &CRI{}
		client, err := cri.client()
		require.NoError(t, err)
		require.IsType(t, &containerdRuntimeClient{}, client)
	})

	t.Run("makeDir fallback uses os.MkdirAll", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "a", "b")
		cri := &CRI{}

		err := cri.makeDir(dir, 0o755)
		require.NoError(t, err)
		info, statErr := os.Stat(dir)
		require.NoError(t, statErr)
		require.True(t, info.IsDir())
	})

	t.Run("afterSignal fallback uses time.After", func(t *testing.T) {
		cri := &CRI{}

		select {
		case <-cri.afterSignal(time.Millisecond):
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout waiting for afterSignal")
		}
	})
}

func TestNewContainerdClient(t *testing.T) {
	t.Run("factory error", func(t *testing.T) {
		wantErr := errors.New("factory failed")
		restore := swapContainerdClientFactory(func(string) (containerdClientAPI, error) {
			return nil, wantErr
		})
		defer restore()

		client, err := newContainerdClient()
		require.Nil(t, client)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("success", func(t *testing.T) {
		fakeAPI := &fakeContainerdClientAPI{}
		restore := swapContainerdClientFactory(func(string) (containerdClientAPI, error) {
			return fakeAPI, nil
		})
		defer restore()

		client, err := newContainerdClient()
		require.NoError(t, err)
		require.IsType(t, &containerdRuntimeClient{}, client)
		require.NoError(t, client.Close())
		require.True(t, fakeAPI.closed)
	})
}

func TestContainerdRuntimeClient(t *testing.T) {
	t.Run("marker", func(t *testing.T) {
		containerdImage{}.isRuntimeImage()
	})

	t.Run("pull success", func(t *testing.T) {
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{}}

		image, err := client.Pull(context.Background(), "docker.io/library/redis:alpine")
		require.NoError(t, err)
		require.IsType(t, containerdImage{}, image)
		image.(containerdImage).isRuntimeImage()
	})

	t.Run("pull error", func(t *testing.T) {
		wantErr := errors.New("pull failed")
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{pullErr: wantErr}}

		image, err := client.Pull(context.Background(), "ref")
		require.Nil(t, image)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("new container wrong image type", func(t *testing.T) {
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{}}

		container, err := client.NewContainer(context.Background(), "c1", fakeImage{})
		require.Nil(t, container)
		require.EqualError(t, err, "unexpected runtime image type cri.fakeImage")
	})

	t.Run("new container error", func(t *testing.T) {
		wantErr := errors.New("new container failed")
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{newContainerErr: wantErr}}

		container, err := client.NewContainer(context.Background(), "c1", containerdImage{})
		require.Nil(t, container)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("new container success", func(t *testing.T) {
		mockContainer := &mockCDContainer{id: "c1"}
		fakeAPI := &fakeContainerdClientAPI{newContainer: mockContainer}
		client := &containerdRuntimeClient{client: fakeAPI}

		container, err := client.NewContainer(context.Background(), "c1", containerdImage{})
		require.NoError(t, err)
		require.IsType(t, &containerdRuntimeContainer{}, container)
		require.Equal(t, "c1", fakeAPI.newContainerID)
		require.Len(t, fakeAPI.newContainerOpts, 3)
	})

	t.Run("containers error", func(t *testing.T) {
		wantErr := errors.New("containers failed")
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{containersErr: wantErr}}

		containers, err := client.Containers(context.Background())
		require.Nil(t, containers)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("containers success", func(t *testing.T) {
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{
			containers: []containerd.Container{&mockCDContainer{id: "c1"}, &mockCDContainer{id: "c2"}},
		}}

		containers, err := client.Containers(context.Background())
		require.NoError(t, err)
		require.Len(t, containers, 2)
		require.IsType(t, &containerdRuntimeContainer{}, containers[0])
	})

	t.Run("list tasks error", func(t *testing.T) {
		wantErr := errors.New("list failed")
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{
			taskService: &fakeContainerdTaskService{listErr: wantErr},
		}}

		tasks, err := client.ListTasks(context.Background())
		require.Nil(t, tasks)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("list tasks success", func(t *testing.T) {
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{
			taskService: &fakeContainerdTaskService{
				resp: &tasksapi.ListTasksResponse{
					Tasks: []*tasktypes.Process{
						{ID: "c1", Pid: 12, Status: tasktypes.Status_RUNNING},
					},
				},
			},
		}}

		tasks, err := client.ListTasks(context.Background())
		require.NoError(t, err)
		require.Equal(t, []runtimeTaskProcess{{ID: "c1", PID: 12, Status: "RUNNING"}}, tasks)
	})

	t.Run("load container error", func(t *testing.T) {
		wantErr := errors.New("load failed")
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{loadContainerErr: wantErr}}

		container, err := client.LoadContainer(context.Background(), "c1")
		require.Nil(t, container)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("load container success", func(t *testing.T) {
		mockContainer := &mockCDContainer{id: "c1"}
		client := &containerdRuntimeClient{client: &fakeContainerdClientAPI{loadContainer: mockContainer}}

		container, err := client.LoadContainer(context.Background(), "c1")
		require.NoError(t, err)
		require.IsType(t, &containerdRuntimeContainer{}, container)
	})
}

func TestContainerdRuntimeContainer(t *testing.T) {
	t.Run("id", func(t *testing.T) {
		container := &containerdRuntimeContainer{container: &mockCDContainer{id: "c1"}}
		require.Equal(t, "c1", container.ID())
	})

	t.Run("info error", func(t *testing.T) {
		wantErr := errors.New("info failed")
		container := &containerdRuntimeContainer{container: &mockCDContainer{infoErr: wantErr}}

		info, err := container.Info(context.Background())
		require.Equal(t, runtimeContainerInfo{}, info)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("info success", func(t *testing.T) {
		container := &containerdRuntimeContainer{container: &mockCDContainer{
			info: containers.Container{
				Image:   "redis:alpine",
				Runtime: containers.RuntimeInfo{Name: "runc"},
			},
		}}

		info, err := container.Info(context.Background())
		require.NoError(t, err)
		require.Equal(t, runtimeContainerInfo{Image: "redis:alpine", Runtime: "runc"}, info)
	})

	t.Run("new task error", func(t *testing.T) {
		wantErr := errors.New("new task failed")
		container := &containerdRuntimeContainer{container: &mockCDContainer{newTaskErr: wantErr}}

		task, err := container.NewTask(context.Background(), "/tmp/example.log")
		require.Nil(t, task)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("new task success", func(t *testing.T) {
		mockTask := &mockCDTask{id: "t1"}
		container := &containerdRuntimeContainer{container: &mockCDContainer{newTask: mockTask}}

		task, err := container.NewTask(context.Background(), "/tmp/example.log")
		require.NoError(t, err)
		require.IsType(t, &containerdRuntimeTask{}, task)
	})

	t.Run("task error", func(t *testing.T) {
		wantErr := errors.New("task failed")
		container := &containerdRuntimeContainer{container: &mockCDContainer{taskErr: wantErr}}

		task, err := container.Task(context.Background())
		require.Nil(t, task)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("task success", func(t *testing.T) {
		mockTask := &mockCDTask{id: "t1"}
		container := &containerdRuntimeContainer{container: &mockCDContainer{task: mockTask}}

		task, err := container.Task(context.Background())
		require.NoError(t, err)
		require.IsType(t, &containerdRuntimeTask{}, task)
	})

	t.Run("delete", func(t *testing.T) {
		mockContainer := &mockCDContainer{}
		container := &containerdRuntimeContainer{container: mockContainer}

		err := container.Delete(context.Background())
		require.NoError(t, err)
		require.True(t, mockContainer.deleteCalled)
	})

	t.Run("stop signal", func(t *testing.T) {
		mockContainer := &mockCDContainer{labels: map[string]string{containerd.StopSignalLabel: "SIGUSR1"}}
		container := &containerdRuntimeContainer{container: mockContainer}

		sig, err := container.StopSignal(context.Background(), syscall.SIGTERM)
		require.NoError(t, err)
		require.Equal(t, syscall.SIGUSR1, sig)
	})
}

func TestContainerdRuntimeTask(t *testing.T) {
	t.Run("id", func(t *testing.T) {
		task := &containerdRuntimeTask{task: &mockCDTask{id: "t1"}}
		require.Equal(t, "t1", task.ID())
	})

	t.Run("wait error", func(t *testing.T) {
		wantErr := errors.New("wait failed")
		task := &containerdRuntimeTask{task: &mockCDTask{waitErr: wantErr}}

		waitC, err := task.Wait(context.Background())
		require.Nil(t, waitC)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("wait closed channel yields zero exit status", func(t *testing.T) {
		ch := make(chan containerd.ExitStatus)
		close(ch)
		task := &containerdRuntimeTask{task: &mockCDTask{waitCh: ch}}

		waitC, err := task.Wait(context.Background())
		require.NoError(t, err)
		exitStatus := <-waitC
		code, exitedAt, resultErr := exitStatus.Result()
		require.NoError(t, resultErr)
		require.Equal(t, uint32(0), code)
		require.True(t, exitedAt.IsZero())
	})

	t.Run("wait success", func(t *testing.T) {
		ch := make(chan containerd.ExitStatus, 1)
		ch <- *containerd.NewExitStatus(7, time.Unix(30, 0), nil)
		task := &containerdRuntimeTask{task: &mockCDTask{waitCh: ch}}

		waitC, err := task.Wait(context.Background())
		require.NoError(t, err)
		code, exitedAt, resultErr := (<-waitC).Result()
		require.NoError(t, resultErr)
		require.Equal(t, uint32(7), code)
		require.Equal(t, time.Unix(30, 0), exitedAt)
	})

	t.Run("start", func(t *testing.T) {
		task := &containerdRuntimeTask{task: &mockCDTask{}}
		require.NoError(t, task.Start(context.Background()))
	})

	t.Run("status error", func(t *testing.T) {
		wantErr := errors.New("status failed")
		task := &containerdRuntimeTask{task: &mockCDTask{statusErr: wantErr}}

		status, err := task.Status(context.Background())
		require.Empty(t, status)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("status success", func(t *testing.T) {
		task := &containerdRuntimeTask{task: &mockCDTask{status: containerd.Status{Status: containerd.Running}}}

		status, err := task.Status(context.Background())
		require.NoError(t, err)
		require.Equal(t, runtimeTaskStatus(containerd.Running), status)
	})

	t.Run("kill", func(t *testing.T) {
		mockTask := &mockCDTask{}
		task := &containerdRuntimeTask{task: mockTask}
		require.NoError(t, task.Kill(context.Background(), syscall.SIGTERM))
		require.Equal(t, []syscall.Signal{syscall.SIGTERM}, mockTask.killed)
	})

	t.Run("delete error", func(t *testing.T) {
		wantErr := errors.New("delete failed")
		task := &containerdRuntimeTask{task: &mockCDTask{deleteErr: wantErr}}

		exitStatus, err := task.Delete(context.Background())
		require.Nil(t, exitStatus)
		require.ErrorIs(t, err, wantErr)
	})

	t.Run("delete nil exit status", func(t *testing.T) {
		task := &containerdRuntimeTask{task: &mockCDTask{}}

		exitStatus, err := task.Delete(context.Background())
		require.NoError(t, err)
		code, exitedAt, resultErr := exitStatus.Result()
		require.NoError(t, resultErr)
		require.Equal(t, uint32(0), code)
		require.True(t, exitedAt.IsZero())
	})

	t.Run("delete success", func(t *testing.T) {
		task := &containerdRuntimeTask{task: &mockCDTask{
			deleteStatus: containerd.NewExitStatus(9, time.Unix(40, 0), nil),
		}}

		exitStatus, err := task.Delete(context.Background())
		require.NoError(t, err)
		code, exitedAt, resultErr := exitStatus.Result()
		require.NoError(t, resultErr)
		require.Equal(t, uint32(9), code)
		require.Equal(t, time.Unix(40, 0), exitedAt)
	})
}

func TestContainerdRuntimeExitStatusResult(t *testing.T) {
	status := containerdRuntimeExitStatus{status: *containerd.NewExitStatus(11, time.Unix(50, 0), nil)}
	code, exitedAt, err := status.Result()
	require.NoError(t, err)
	require.Equal(t, uint32(11), code)
	require.Equal(t, time.Unix(50, 0), exitedAt)
}

func TestCRI_StopAndDeleteTaskBridge(t *testing.T) {
	waitCh := make(chan containerd.ExitStatus, 1)
	waitCh <- *containerd.NewExitStatus(0, time.Unix(60, 0), nil)

	task := &mockCDTask{
		id:           "t1",
		status:       containerd.Status{Status: containerd.Running},
		waitCh:       waitCh,
		deleteStatus: containerd.NewExitStatus(0, time.Unix(61, 0), nil),
	}
	container := &mockCDContainer{
		id:     "c1",
		labels: map[string]string{containerd.StopSignalLabel: "SIGUSR1"},
	}
	cri := newTestCRI(nil)

	err := cri.StopAndDeleteTask(context.Background(), "example", container, task, time.Second)
	require.NoError(t, err)
	require.Equal(t, []syscall.Signal{syscall.SIGUSR1}, task.killed)
}

func swapContainerdClientFactory(factory func(string) (containerdClientAPI, error)) func() {
	previous := containerdClientFactory
	containerdClientFactory = factory
	return func() {
		containerdClientFactory = previous
	}
}

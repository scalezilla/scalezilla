package cri

import (
	"context"
	"syscall"

	tasksapi "github.com/containerd/containerd/api/services/tasks/v1"
	containerdtypes "github.com/containerd/containerd/api/types"
	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/core/containers"
	"github.com/containerd/containerd/v2/pkg/cio"
	"github.com/containerd/containerd/v2/pkg/oci"
	"github.com/containerd/typeurl/v2"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeContainerdClientAPI struct {
	pullErr          error
	newContainerErr  error
	containersErr    error
	loadContainerErr error
	taskService      tasksapi.TasksClient
	containers       []containerd.Container
	newContainer     containerd.Container
	loadContainer    containerd.Container
	newContainerID   string
	newContainerOpts []containerd.NewContainerOpts
	closed           bool
}

func (f *fakeContainerdClientAPI) Pull(ctx context.Context, ref string, opts ...containerd.RemoteOpt) (containerd.Image, error) {
	return nil, f.pullErr
}

func (f *fakeContainerdClientAPI) NewContainer(ctx context.Context, id string, opts ...containerd.NewContainerOpts) (containerd.Container, error) {
	f.newContainerID = id
	f.newContainerOpts = opts
	return f.newContainer, f.newContainerErr
}

func (f *fakeContainerdClientAPI) Containers(ctx context.Context, filters ...string) ([]containerd.Container, error) {
	return f.containers, f.containersErr
}

func (f *fakeContainerdClientAPI) TaskService() tasksapi.TasksClient {
	if f.taskService == nil {
		return &fakeContainerdTaskService{}
	}
	return f.taskService
}

func (f *fakeContainerdClientAPI) LoadContainer(ctx context.Context, id string) (containerd.Container, error) {
	return f.loadContainer, f.loadContainerErr
}

func (f *fakeContainerdClientAPI) Close() error {
	f.closed = true
	return nil
}

type fakeContainerdTaskService struct {
	resp    *tasksapi.ListTasksResponse
	listErr error
}

func (f *fakeContainerdTaskService) Create(ctx context.Context, in *tasksapi.CreateTaskRequest, opts ...grpc.CallOption) (*tasksapi.CreateTaskResponse, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Start(ctx context.Context, in *tasksapi.StartRequest, opts ...grpc.CallOption) (*tasksapi.StartResponse, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Delete(ctx context.Context, in *tasksapi.DeleteTaskRequest, opts ...grpc.CallOption) (*tasksapi.DeleteResponse, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) DeleteProcess(ctx context.Context, in *tasksapi.DeleteProcessRequest, opts ...grpc.CallOption) (*tasksapi.DeleteResponse, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Get(ctx context.Context, in *tasksapi.GetRequest, opts ...grpc.CallOption) (*tasksapi.GetResponse, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) List(ctx context.Context, in *tasksapi.ListTasksRequest, opts ...grpc.CallOption) (*tasksapi.ListTasksResponse, error) {
	if f.resp == nil {
		f.resp = &tasksapi.ListTasksResponse{}
	}
	return f.resp, f.listErr
}

func (f *fakeContainerdTaskService) Kill(ctx context.Context, in *tasksapi.KillRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Exec(ctx context.Context, in *tasksapi.ExecProcessRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) ResizePty(ctx context.Context, in *tasksapi.ResizePtyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) CloseIO(ctx context.Context, in *tasksapi.CloseIORequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Pause(ctx context.Context, in *tasksapi.PauseTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Resume(ctx context.Context, in *tasksapi.ResumeTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) ListPids(ctx context.Context, in *tasksapi.ListPidsRequest, opts ...grpc.CallOption) (*tasksapi.ListPidsResponse, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Checkpoint(ctx context.Context, in *tasksapi.CheckpointTaskRequest, opts ...grpc.CallOption) (*tasksapi.CheckpointTaskResponse, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Update(ctx context.Context, in *tasksapi.UpdateTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Metrics(ctx context.Context, in *tasksapi.MetricsRequest, opts ...grpc.CallOption) (*tasksapi.MetricsResponse, error) {
	return nil, nil
}

func (f *fakeContainerdTaskService) Wait(ctx context.Context, in *tasksapi.WaitRequest, opts ...grpc.CallOption) (*tasksapi.WaitResponse, error) {
	return nil, nil
}

type mockCDContainer struct {
	id           string
	info         containers.Container
	infoErr      error
	newTask      containerd.Task
	newTaskErr   error
	task         containerd.Task
	taskErr      error
	deleteErr    error
	deleteCalled bool
	labels       map[string]string
	labelsErr    error
}

func (m *mockCDContainer) ID() string {
	return m.id
}

func (m *mockCDContainer) Info(ctx context.Context, opts ...containerd.InfoOpts) (containers.Container, error) {
	return m.info, m.infoErr
}

func (m *mockCDContainer) Delete(ctx context.Context, opts ...containerd.DeleteOpts) error {
	m.deleteCalled = true
	return m.deleteErr
}

func (m *mockCDContainer) NewTask(ctx context.Context, creator cio.Creator, opts ...containerd.NewTaskOpts) (containerd.Task, error) {
	return m.newTask, m.newTaskErr
}

func (m *mockCDContainer) Spec(ctx context.Context) (*oci.Spec, error) {
	return &oci.Spec{}, nil
}

func (m *mockCDContainer) Task(ctx context.Context, attach cio.Attach) (containerd.Task, error) {
	return m.task, m.taskErr
}

func (m *mockCDContainer) Image(ctx context.Context) (containerd.Image, error) {
	return nil, nil
}

func (m *mockCDContainer) Labels(ctx context.Context) (map[string]string, error) {
	return m.labels, m.labelsErr
}

func (m *mockCDContainer) SetLabels(ctx context.Context, labels map[string]string) (map[string]string, error) {
	m.labels = labels
	return labels, nil
}

func (m *mockCDContainer) Extensions(ctx context.Context) (map[string]typeurl.Any, error) {
	return nil, nil
}

func (m *mockCDContainer) Update(ctx context.Context, opts ...containerd.UpdateContainerOpts) error {
	return nil
}

func (m *mockCDContainer) Checkpoint(ctx context.Context, ref string, opts ...containerd.CheckpointOpts) (containerd.Image, error) {
	return nil, nil
}

func (m *mockCDContainer) Restore(ctx context.Context, creator cio.Creator, ref string) (int, error) {
	return 0, nil
}

type mockCDTask struct {
	id           string
	pid          uint32
	status       containerd.Status
	statusErr    error
	waitCh       chan containerd.ExitStatus
	waitErr      error
	startErr     error
	deleteStatus *containerd.ExitStatus
	deleteErr    error
	killErr      error
	killed       []syscall.Signal
}

func (m *mockCDTask) ID() string {
	return m.id
}

func (m *mockCDTask) Pid() uint32 {
	return m.pid
}

func (m *mockCDTask) Start(ctx context.Context) error {
	return m.startErr
}

func (m *mockCDTask) Delete(ctx context.Context, opts ...containerd.ProcessDeleteOpts) (*containerd.ExitStatus, error) {
	return m.deleteStatus, m.deleteErr
}

func (m *mockCDTask) Kill(ctx context.Context, signal syscall.Signal, opts ...containerd.KillOpts) error {
	m.killed = append(m.killed, signal)
	return m.killErr
}

func (m *mockCDTask) Wait(ctx context.Context) (<-chan containerd.ExitStatus, error) {
	return m.waitCh, m.waitErr
}

func (m *mockCDTask) CloseIO(ctx context.Context, opts ...containerd.IOCloserOpts) error {
	return nil
}

func (m *mockCDTask) Resize(ctx context.Context, w, h uint32) error {
	return nil
}

func (m *mockCDTask) IO() cio.IO {
	return nil
}

func (m *mockCDTask) Status(ctx context.Context) (containerd.Status, error) {
	return m.status, m.statusErr
}

func (m *mockCDTask) Pause(ctx context.Context) error {
	return nil
}

func (m *mockCDTask) Resume(ctx context.Context) error {
	return nil
}

func (m *mockCDTask) Exec(ctx context.Context, id string, spec *specs.Process, creator cio.Creator) (containerd.Process, error) {
	return nil, nil
}

func (m *mockCDTask) Pids(ctx context.Context) ([]containerd.ProcessInfo, error) {
	return nil, nil
}

func (m *mockCDTask) Checkpoint(ctx context.Context, opts ...containerd.CheckpointTaskOpts) (containerd.Image, error) {
	return nil, nil
}

func (m *mockCDTask) Update(ctx context.Context, opts ...containerd.UpdateTaskOpts) error {
	return nil
}

func (m *mockCDTask) LoadProcess(ctx context.Context, id string, attach cio.Attach) (containerd.Process, error) {
	return nil, nil
}

func (m *mockCDTask) Metrics(ctx context.Context) (*containerdtypes.Metric, error) {
	return nil, nil
}

func (m *mockCDTask) Spec(ctx context.Context) (*oci.Spec, error) {
	return &oci.Spec{}, nil
}

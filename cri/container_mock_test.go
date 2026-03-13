package cri

import (
	"context"
	"syscall"
	"time"
)

type fakeClient struct {
	pullImage        runtimeImage
	pullErr          error
	newContainer     runtimeContainer
	newContainerErr  error
	containers       []runtimeContainer
	containersErr    error
	taskList         []runtimeTaskProcess
	listTasksErr     error
	loadContainer    runtimeContainer
	loadContainerErr error
	pullRef          string
	newContainerID   string
	newContainerImg  runtimeImage
	closed           bool
}

func (f *fakeClient) Pull(ctx context.Context, ref string) (runtimeImage, error) {
	f.pullRef = ref
	return f.pullImage, f.pullErr
}

func (f *fakeClient) NewContainer(ctx context.Context, id string, image runtimeImage) (runtimeContainer, error) {
	f.newContainerID = id
	f.newContainerImg = image
	return f.newContainer, f.newContainerErr
}

func (f *fakeClient) Containers(ctx context.Context) ([]runtimeContainer, error) {
	return f.containers, f.containersErr
}

func (f *fakeClient) ListTasks(ctx context.Context) ([]runtimeTaskProcess, error) {
	return f.taskList, f.listTasksErr
}

func (f *fakeClient) LoadContainer(ctx context.Context, id string) (runtimeContainer, error) {
	return f.loadContainer, f.loadContainerErr
}

func (f *fakeClient) Close() error {
	f.closed = true
	return nil
}

type fakeImage struct{}

func (fakeImage) isRuntimeImage() {}

type fakeContainer struct {
	id             string
	info           runtimeContainerInfo
	infoErr        error
	newTask        runtimeTask
	newTaskErr     error
	task           runtimeTask
	taskErr        error
	deleteErr      error
	stopSignal     syscall.Signal
	stopSignalErr  error
	newTaskLogFile string
	deleted        bool
}

func (f *fakeContainer) ID() string {
	return f.id
}

func (f *fakeContainer) Info(ctx context.Context) (runtimeContainerInfo, error) {
	return f.info, f.infoErr
}

func (f *fakeContainer) NewTask(ctx context.Context, logFile string) (runtimeTask, error) {
	f.newTaskLogFile = logFile
	return f.newTask, f.newTaskErr
}

func (f *fakeContainer) Task(ctx context.Context) (runtimeTask, error) {
	return f.task, f.taskErr
}

func (f *fakeContainer) Delete(ctx context.Context) error {
	f.deleted = true
	return f.deleteErr
}

func (f *fakeContainer) StopSignal(ctx context.Context, defaultSignal syscall.Signal) (syscall.Signal, error) {
	if f.stopSignalErr != nil {
		return 0, f.stopSignalErr
	}
	if f.stopSignal == 0 {
		return defaultSignal, nil
	}
	return f.stopSignal, nil
}

type fakeTask struct {
	id               string
	waitCh           chan runtimeExitStatus
	waitErr          error
	startErr         error
	status           runtimeTaskStatus
	statusErr        error
	killErrs         []error
	deleteExitStatus runtimeExitStatus
	deleteErr        error
	killedSignals    []syscall.Signal
	onKill           func(syscall.Signal)
	started          bool
	deleted          bool
}

func (f *fakeTask) ID() string {
	return f.id
}

func (f *fakeTask) Wait(ctx context.Context) (<-chan runtimeExitStatus, error) {
	return f.waitCh, f.waitErr
}

func (f *fakeTask) Start(ctx context.Context) error {
	f.started = true
	return f.startErr
}

func (f *fakeTask) Status(ctx context.Context) (runtimeTaskStatus, error) {
	return f.status, f.statusErr
}

func (f *fakeTask) Kill(ctx context.Context, sig syscall.Signal) error {
	f.killedSignals = append(f.killedSignals, sig)
	if f.onKill != nil {
		f.onKill(sig)
	}
	if len(f.killErrs) == 0 {
		return nil
	}
	err := f.killErrs[0]
	f.killErrs = f.killErrs[1:]
	return err
}

func (f *fakeTask) Delete(ctx context.Context) (runtimeExitStatus, error) {
	f.deleted = true
	return f.deleteExitStatus, f.deleteErr
}

type fakeExitStatus struct {
	code uint32
	at   time.Time
	err  error
}

func (f fakeExitStatus) Result() (uint32, time.Time, error) {
	return f.code, f.at, f.err
}

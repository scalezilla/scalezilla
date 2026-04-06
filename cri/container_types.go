package cri

import (
	"context"
	"os"
	"syscall"
	"time"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/rs/zerolog"
)

const (
	// containerdAddress is containerd socket address
	containerdAddress = "/run/containerd/containerd.sock"

	// constainerdNamespace is the containerd namespace used by scalezilla
	constainerdNamespace = "scalezilla"
)

// CRI is used by NewCRI
type CRI struct {
	log       *zerolog.Logger
	newClient func() (runtimeClient, error)
	mkdirAll  func(path string, perm os.FileMode) error
	after     func(time.Duration) <-chan time.Time
}

// ImageSpec holds container image specs
type ImageSpec struct {
	Image string
}

// CreateContainerSpec holds container specs to create a container
type CreateContainerSpec struct {
	Namespace      string
	ContainerID    string
	Labels         map[string]string
	Image          ImageSpec
	DefaultLogPath string
}

type DeploymentSpec struct {
	Deployment DeploymentConfigSpec `hcl:"deployment,block"`
}

type DeploymentConfigSpec struct {
	Name      string            `hcl:"deployment,label"`
	Kind      string            `hcl:"kind,optional"`
	Namespace string            `hcl:"namespace,optional"`
	Metadata  map[string]string `hcl:"metadata,optional"`
	Pod       DeploymentPodSpec `hcl:"pod,block"`
}

type DeploymentPodSpec struct {
	Name      string                  `hcl:"pod,label"`
	Container DeploymentContainerSpec `hcl:"container,block"`
}

type DeploymentContainerSpec struct {
	Name      string         `hcl:"container,label"`
	Image     string         `hcl:"image"`
	Resources *ResourcesSpec `hcl:"resources,block"`
}

type ResourcesSpec struct {
	CPU    uint64 `hcl:"cpu"`
	Memory uint64 `hcl:"memory"`
}

// ContainerList returns the container with its status
type ContainerList struct {
	Namespace string
	ID        string
	PID       uint32
	Image     string
	Runtime   string
	Status    string
	CreatedAt time.Time
}

// ContainerRuntime is an interface implements containers requirements
type ContainerRuntime interface {
	CreateContainer(ctx context.Context, spec CreateContainerSpec) error
	ListContainer(ctx context.Context, namespace string) ([]ContainerList, error)
	DeleteContainer(ctx context.Context, namespace, containerID string, stopTimeout time.Duration) error
	StopAndDeleteTask(ctx context.Context, namespace string, container containerd.Container, task containerd.Task, stopTimeout time.Duration) error
}

// runtimeClient abstracts the client operations used by CRI.
// It keeps unit tests decoupled from the concrete containerd client.
type runtimeClient interface {
	Pull(ctx context.Context, ref string) (runtimeImage, error)
	NewContainer(ctx context.Context, id string, image runtimeImage, labels, additionalContainerLabels map[string]string) (runtimeContainer, error)
	Containers(ctx context.Context) ([]runtimeContainer, error)
	ListTasks(ctx context.Context) ([]runtimeTaskProcess, error)
	LoadContainer(ctx context.Context, id string) (runtimeContainer, error)
	Close() error
}

// runtimeImage is a marker for images returned by runtimeClient.
// It prevents tests from depending on containerd image types directly.
type runtimeImage interface {
	isRuntimeImage()
}

// runtimeContainer exposes the container operations used by CRI.
// The interface is intentionally narrower than containerd.Container.
type runtimeContainer interface {
	ID() string
	Info(ctx context.Context) (runtimeContainerInfo, error)
	NewTask(ctx context.Context, logFile string) (runtimeTask, error)
	Task(ctx context.Context) (runtimeTask, error)
	Delete(ctx context.Context) error
	StopSignal(ctx context.Context, defaultSignal syscall.Signal) (syscall.Signal, error)
}

// runtimeTask exposes the task operations used during container lifecycle.
// It allows task shutdown logic to be unit-tested with fakes.
type runtimeTask interface {
	ID() string
	Wait(ctx context.Context) (<-chan runtimeExitStatus, error)
	Start(ctx context.Context) error
	Status(ctx context.Context) (runtimeTaskStatus, error)
	Kill(ctx context.Context, sig syscall.Signal) error
	Delete(ctx context.Context) (runtimeExitStatus, error)
}

// runtimeContainerInfo holds the container fields CRI needs to render.
type runtimeContainerInfo struct {
	Image     string
	Runtime   string
	Labels    map[string]string
	CreatedAt time.Time
}

// runtimeTaskProcess is the reduced task view used by ListContainer.
type runtimeTaskProcess struct {
	ID     string
	PID    uint32
	Status string
}

// runtimeTaskStatus mirrors the task status values consumed by CRI.
type runtimeTaskStatus string

const runtimeTaskStopped runtimeTaskStatus = "stopped"

// runtimeExitStatus exposes the exit information read from wait and delete calls.
type runtimeExitStatus interface {
	Result() (uint32, time.Time, error)
}

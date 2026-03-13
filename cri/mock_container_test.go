package cri

import (
	"context"
	"time"

	containerd "github.com/containerd/containerd/v2/client"
)

type mockContainer struct {
	called                                   bool
	errCreateContainer, errListContainer     error
	errDeleteContainer, errStopAndDeleteTask error
}

func (m *mockContainer) CreateContainer(ctx context.Context, spec CreateContainerSpec) error {
	m.called = true
	return m.errCreateContainer
}

func (m *mockContainer) ListContainer(ctx context.Context, namespace string) ([]ContainerList, error) {
	m.called = true
	return nil, m.errListContainer
}

func (m *mockContainer) DeleteContainer(ctx context.Context, namespace, containerID string, stopTimeout time.Duration) error {
	m.called = true
	return m.errDeleteContainer
}

func (m *mockContainer) StopAndDeleteTask(ctx context.Context, namespace string, container containerd.Container, task containerd.Task, stopTimeout time.Duration) error {
	m.called = true
	return m.errStopAndDeleteTask
}

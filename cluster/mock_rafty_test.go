package cluster

import (
	"time"

	"github.com/Lord-Y/rafty"
)

type mockRafty struct {
	called                         bool
	err, errBootstrap              error
	bootstrapped, isLeader, leader bool
	leaderAddress, leaderId        string
	raftyStatus                    rafty.Status
}

func (m *mockRafty) Start() error {
	m.called = true
	return m.err
}

func (m *mockRafty) Stop() {
	m.called = true
}

func (m *mockRafty) IsBootstrapped() bool {
	m.called = true
	return m.bootstrapped
}

func (m *mockRafty) SubmitCommand(timeout time.Duration, logKind rafty.LogKind, command []byte) ([]byte, error) {
	m.called = true
	return nil, m.err
}

func (m *mockRafty) BootstrapCluster(timeout time.Duration) error {
	m.called = true
	return m.errBootstrap
}

func (m *mockRafty) IsLeader() bool {
	m.called = true
	return m.isLeader
}

func (m *mockRafty) Status() rafty.Status {
	m.called = true
	return m.raftyStatus
}

func (m *mockRafty) Leader() (bool, string, string) {
	m.called = true
	return m.leader, m.leaderAddress, m.leaderId
}

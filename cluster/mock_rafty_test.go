package cluster

import (
	"time"

	"github.com/Lord-Y/rafty"
)

type mockRafty struct {
	called       bool
	err          error
	bootstrapped bool
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

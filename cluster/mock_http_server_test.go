package cluster

import "context"

type mockHTTPServer struct {
	called bool
	err    error
}

func (m *mockHTTPServer) ListenAndServe() error {
	m.called = true
	return m.err
}

func (m *mockHTTPServer) Shutdown(ctx context.Context) error {
	m.called = true
	return m.err
}

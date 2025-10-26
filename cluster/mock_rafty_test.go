package cluster

type mockRafty struct {
	called bool
	err    error
}

func (m *mockRafty) Start() error {
	m.called = true
	return m.err
}

func (m *mockRafty) Stop() {
	m.called = true
}

package cluster

import (
	"github.com/Lord-Y/rafty"
)

// deploymentSet will add key/value to the deployment store.
// An error will be returned if any
func (m *memoryStore) deploymentSet(log *rafty.LogEntry, config deploymentState) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// delete existing key if exist
	// That will allows us to cleanly perform snapshots
	// when required by removed overriden keys and reduce
	// disk space and amount of time to restore data
	if _, ok := m.deployment[config.Name]; ok {
		delete(m.logs, m.deployment[config.Name].index)
	}

	m.logs[log.Index] = log
	config.index = log.Index
	m.deployment[config.Name] = config
	return nil
}

// deploymentGet will fetch provided key from the deployment store.
// An error will be returned if the key is not found
func (m *memoryStore) deploymentGet(key []byte) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keyName := string(key)
	if _, ok := m.deployment[keyName]; ok {
		return m.logs[m.deployment[keyName].index].Command, nil
	}
	return nil, rafty.ErrKeyNotFound
}

// deploymentExist will return true if the deployment exist
func (m *memoryStore) deploymentExist(key []byte) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.deployment[string(key)]; ok {
		return true
	}
	return false
}

// deploymentDelete will delete provided key from the deployment store.
// An error will be returned if the key is not found
func (m *memoryStore) deploymentDelete(key []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data := string(key)
	if _, ok := m.deployment[data]; ok {
		delete(m.logs, m.deployment[data].index)
		delete(m.deployment, string(key))
	}
}

// deploymentGetAll will fetch all deployments from the deployment store.
// An error will be returned if the any
func (m *memoryStore) deploymentGetAll() (z []*deploymentState, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.deployment) == 0 {
		return nil, nil
	}

	for _, v := range m.deployment {
		z = append(z, &v)
	}
	return
}

// deploymentEncoded will fetch all token from the deployment store.
// deployment will be binary encoded when command is forwarded to the leader.
// An error will be returned if the any
func (m *memoryStore) deploymentEncoded(cmd deploymentState) (u []byte, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.deployment) == 0 {
		return nil, nil
	}

	if cmd.Kind == deploymentCommandGet {
		value, err := m.deploymentGet([]byte(cmd.Name))
		if err != nil {
			return nil, err
		}
		return value, nil
	}

	for _, v := range m.deployment {
		u = append(u, m.logs[v.index].Command...)
	}
	return
}

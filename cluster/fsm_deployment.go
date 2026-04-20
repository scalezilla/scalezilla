package cluster

import (
	"fmt"

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
	name := fmt.Sprintf("%s-%s", config.Namespace, config.Name)
	if _, ok := m.deployment[name]; ok {
		delete(m.logs, m.deployment[name].index)
	}

	m.logs[log.Index] = log
	config.index = log.Index
	m.deployment[name] = config
	return nil
}

// deploymentGet will fetch provided key from the deployment store.
// An error will be returned if the key is not found
func (m *memoryStore) deploymentGet(namespace, deploymentName []byte) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ns := string(namespace)
	dname := string(deploymentName)
	name := fmt.Sprintf("%s-%s", ns, dname)
	if _, ok := m.deployment[name]; ok {
		return m.logs[m.deployment[name].index].Command, nil
	}
	return nil, rafty.ErrKeyNotFound
}

// deploymentExist will return true if the deployment exist
func (m *memoryStore) deploymentExist(namespace, deploymentName []byte) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.deployment[fmt.Sprintf("%s-%s", string(namespace), string(deploymentName))]; ok {
		return true
	}
	return false
}

// deploymentDelete will delete provided key from the deployment store.
// An error will be returned if the key is not found
func (m *memoryStore) deploymentDelete(namespace, deploymentName []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ns := string(namespace)
	dname := string(deploymentName)
	data := fmt.Sprintf("%s-%s", ns, dname)
	if _, ok := m.deployment[data]; ok {
		delete(m.logs, m.deployment[data].index)
		delete(m.deployment, data)
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
		value, err := m.deploymentGet([]byte(cmd.Namespace), []byte(cmd.Name))
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

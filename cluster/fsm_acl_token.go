package cluster

import (
	"github.com/Lord-Y/rafty"
)

// aclTokenSet will add key/value to the aclToken store.
// An error will be returned if any
func (m *memoryStore) aclTokenSet(log *rafty.LogEntry, token aclTokenCommand) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// delete existing key if exist
	// That will allows us to cleanly perform snapshots
	// when required by removed overriden keys and reduce
	// disk space and amount of time to restore data
	if _, ok := m.aclToken[token.AccessorID]; ok {
		delete(m.logs, m.aclToken[token.AccessorID].index)
	}

	m.logs[log.Index] = log
	m.aclToken[token.AccessorID] = dataACLToken{index: log.Index, value: token}
	return nil
}

// aclTokenGet will fetch provided key from the aclToken store.
// An error will be returned if the key is not found
func (m *memoryStore) aclTokenGet(key []byte) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keyName := string(key)
	if _, ok := m.aclToken[keyName]; ok {
		return m.logs[m.aclToken[keyName].index].Command, nil
	}
	return nil, rafty.ErrKeyNotFound
}

// aclTokenExist will return true if the tken exist
func (m *memoryStore) aclTokenExist(key []byte) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.aclToken[string(key)]; ok {
		return true
	}
	return false
}

// aclTokenDelete will delete provided key from the aclToken store.
// An error will be returned if the key is not found
func (m *memoryStore) aclTokenDelete(key []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data := string(key)
	if _, ok := m.aclToken[data]; ok {
		delete(m.logs, m.aclToken[data].index)
		delete(m.aclToken, string(key))
	}
}

// usersGetAll will fetch all users from the users store.
// An error will be returned if the any
func (m *memoryStore) aclTokenGetAll() (z []*AclToken, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.aclToken) == 0 {
		return nil, nil
	}

	for k, v := range m.aclToken {
		z = append(z, &AclToken{AccessorID: k,
			Token:        v.value.Token,
			InitialToken: v.value.InitialToken},
		)
	}
	return
}

// aclTokenEncoded will fetch all token from the aclToken store.
// tokens will be binary encoded when command is forwarded
// to the leader.
// An error will be returned if the any
func (m *memoryStore) aclTokenEncoded(cmd aclTokenCommand) (u []byte, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.aclToken) == 0 {
		return nil, nil
	}

	if cmd.Kind == aclTokenCommandGet {
		value, err := m.aclTokenGet([]byte(cmd.AccessorID))
		if err != nil {
			return nil, err
		}
		// TODO: I need to get why I was doing the following
		// because in the for loop below I wasn't
		// token, err := aclTokenUnmarshal(value)
		// if err != nil {
		// 	return nil, err
		// }
		// fmt.Printf("TOKEN %+v", token)

		// return json.Marshal(AclToken{
		// 	Firstname: cmd.Key,
		// 	Lastname:  string(value),
		// })

		return value, nil
	}

	// var tokens []AclToken
	for _, v := range m.aclToken {
		// data := AclToken{
		// 	Firstname: string(k),
		// 	Lastname:  string(v.value),
		// }
		// token, err := aclTokenUnMarshal(v.value)
		// if err != nil {
		// 	return nil, err
		// }
		// tokens = append(tokens, token)
		u = append(u, m.logs[v.index].Command...)
	}
	// return json.Marshal(tokens)
	return
}

package cluster

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/Lord-Y/rafty"
	"github.com/stretchr/testify/assert"
)

func TestCluster_fsm_deployment(t *testing.T) {
	assert := assert.New(t)

	t.Run("set", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		cmd := deploymentState{
			Kind:      deploymentCommandGet,
			Namespace: "default",
			Name:      "redis",
		}
		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now(),
			ReplicaSetID: "abcd1234",
		}
		cmd2 := deploymentState{
			Kind:               deploymentCommandSet,
			Namespace:          "default",
			Name:               "redis",
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		assert.Nil(cluster.fsm.memoryStore.deploymentSet(entry, cmd2))
		assert.Nil(cluster.fsm.memoryStore.deploymentSet(entry, cmd2)) // for delete, keep it
	})

	t.Run("get", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		namespace := "default"
		name := "redis"
		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind:               deploymentCommandSet,
			Namespace:          namespace,
			Name:               name,
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		r, err := cluster.fsm.memoryStore.deploymentGet([]byte(namespace), []byte(name))
		assert.Nil(r)
		assert.ErrorIs(err, rafty.ErrKeyNotFound)
		assert.Nil(cluster.fsm.memoryStore.deploymentSet(entry, cmd))
		r, err = cluster.fsm.memoryStore.deploymentGet([]byte(namespace), []byte(name))
		assert.NotNil(r)
		assert.Nil(err)
	})

	t.Run("exist", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		namespace := "default"
		name := "redis"
		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind:               deploymentCommandSet,
			Namespace:          namespace,
			Name:               name,
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		assert.Equal(false, cluster.fsm.memoryStore.deploymentExist([]byte(namespace), []byte(name)))
		assert.Nil(cluster.fsm.memoryStore.deploymentSet(entry, cmd))
		assert.Equal(true, cluster.fsm.memoryStore.deploymentExist([]byte(namespace), []byte(name)))
	})

	t.Run("delete", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		namespace := "default"
		name := "redis"
		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind:               deploymentCommandSet,
			Namespace:          namespace,
			Name:               name,
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		cluster.fsm.memoryStore.deploymentDelete([]byte(namespace), []byte(name))
		assert.Nil(cluster.fsm.memoryStore.deploymentSet(entry, cmd))
		cluster.fsm.memoryStore.deploymentDelete([]byte(namespace), []byte(name))
	})

	t.Run("get_all", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		namespace := "default"
		name := "redis"
		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		cmd := deploymentState{
			Kind:               deploymentCommandSet,
			Namespace:          namespace,
			Name:               name,
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(cmd, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Command: buffer.Bytes(),
		}

		r, err := cluster.fsm.memoryStore.deploymentGetAll()
		assert.Nil(r)
		assert.Nil(err)
		assert.Nil(cluster.fsm.memoryStore.deploymentSet(entry, cmd))
		r, err = cluster.fsm.memoryStore.deploymentGetAll()
		assert.NotNil(r)
		assert.Nil(err)
	})

	t.Run("encoded", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

		namespace := "default"
		name := "redis"
		content := make(map[uint64]deploymentContent)
		key := uint64(0)
		content[key] = deploymentContent{
			IsStable:     true,
			RawContent:   "aaa",
			Version:      1,
			CreatedAt:    time.Now().Round(0),
			ReplicaSetID: "abcd1234",
		}
		set := deploymentState{
			Kind:               deploymentCommandSet,
			Namespace:          namespace,
			Name:               name,
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		get := deploymentState{
			Kind:      deploymentCommandGet,
			Namespace: namespace,
			Name:      name,
		}
		fakeGet := deploymentState{
			Kind:      deploymentCommandGet,
			Namespace: namespace,
			Name:      name + "aa",
		}
		getAll := deploymentState{
			Kind: deploymentCommandGetAll,
		}

		buffer := new(bytes.Buffer)
		err = deploymentEncodeCommand(set, buffer)
		assert.Nil(err)
		entry := &rafty.LogEntry{
			LogType: uint32(rafty.LogReplication),
			Term:    1,
			Index:   1,
			Command: buffer.Bytes(),
		}

		assert.Nil(cluster.fsm.memoryStore.deploymentSet(entry, set))
		r, err := cluster.fsm.memoryStore.deploymentEncoded(get)
		assert.NotNil(r)
		assert.Nil(err)

		r, err = cluster.fsm.memoryStore.deploymentEncoded(fakeGet)
		assert.Nil(r)
		assert.NotNil(err, rafty.ErrKeyNotFound)

		r, err = cluster.fsm.memoryStore.deploymentEncoded(getAll)
		assert.NotNil(r)
		assert.Nil(err)
	})
}

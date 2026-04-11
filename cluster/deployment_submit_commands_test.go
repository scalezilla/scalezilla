package cluster

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCLuster_deployment_submit_commands(t *testing.T) {
	assert := assert.New(t)

	t.Run("write_success", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

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
			Name:               "redis",
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		mock := mockRafty{}
		cluster.rafty = &mock

		assert.Nil(cluster.submitCommandDeploymentWrite(time.Second, cmd))
		assert.Equal(true, mock.called)
	})

	t.Run("write_error", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: false, dev: true}
		cluster := makeBasicCluster(cfg)
		store, err := cluster.buildStore()
		defer func() {
			_ = os.RemoveAll(cluster.config.DataDir)
		}()
		assert.Nil(err)
		cluster.fsm = newFSM(store)

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
			Name:               "redis",
			NewRollingVersion:  -1,
			CurrentUsedVersion: 1,
			Content:            content,
		}

		cluster.di.deploymentEncodeCommandFunc = func(cmd deploymentState, w io.Writer) error {
			return errors.New("deployment encode commnd error")
		}
		assert.Error(cluster.submitCommandDeploymentWrite(time.Second, cmd))

		mock := mockRafty{
			err: errors.New("rafty submit command error"),
		}
		cluster.rafty = &mock

		cluster.di.deploymentEncodeCommandFunc = deploymentEncodeCommand

		assert.Error(cluster.submitCommandDeploymentWrite(time.Second, cmd))
		assert.Equal(true, mock.called)
	})
}

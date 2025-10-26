package cluster

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAPIServer(t *testing.T) {
	assert := assert.New(t)

	t.Run("start_success", func(t *testing.T) {
		cluster := makeBasicCluster(true)
		cluster.newAPIServer()

		go func() {
			assert.Nil(cluster.startAPIServer())
		}()
		time.Sleep(time.Second)
		assert.Nil(cluster.stopAPIServer())
	})

	t.Run("start_failure", func(t *testing.T) {
		mock := mockHTTPServer{
			err: errors.New("shutdown error"),
		}
		cluster := &Cluster{
			apiServer: &mock,
		}

		assert.Error(cluster.startAPIServer())
		assert.Equal(true, mock.called)
	})

	t.Run("stop_failure", func(t *testing.T) {
		mock := mockHTTPServer{
			err: errors.New("shutdown error"),
		}
		cluster := &Cluster{
			apiServer: &mock,
		}

		assert.Error(cluster.stopAPIServer())
		assert.Equal(true, mock.called)
	})
}

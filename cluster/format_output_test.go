package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_print_table_nodes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body := []byte(`[{"id":"127.0.0.1:20001","name":"127.0.0.1:20001","address":"127.0.0.1:20001","kind":"server","leader":true,"nodePool":"default"},{"id":"127.0.0.1:20005","name":"fake","address":"127.0.0.1:20005","kind":"client","leader":false,"nodePool":"default"},{"id":"127.0.0.1:20004","name":"127.0.0.1:20004","address":"127.0.0.1:20004","kind":"server","leader":false,"nodePool":"default"},{"id":"127.0.0.1:20007","name":"127.0.0.1:20007","address":"127.0.0.1:20007","kind":"server","leader":false,"nodePool":"default"}]`)

		printTableNodesList(body)
	})

	t.Run("error", func(t *testing.T) {
		body := []byte(`[{"id":"127.0.0.1:20001"`)

		printTableNodesList(body)
	})
}

func TestCluster_decode_error(t *testing.T) {
	assert := assert.New(t)

	t.Run("success", func(t *testing.T) {
		body := []byte(`{"error":"cluster not boostrapped"}`)
		assert.ErrorContains(decodeError(body), ErrClusterNotBootstrapped.Error())
	})

	t.Run("error", func(t *testing.T) {
		body := []byte(`{"errorr":`)
		assert.ErrorContains(decodeError(body), "unexpected end of JSON input")
	})
}

func TestCluster_print_table_pods(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body := []byte(`[{"namespace":"default","id":"nginx-test-nginx-container","image":"docker.io/library/nginx:latest","pid":5086,"runtime":"io.containerd.runc.v2","status":"RUNNING","created_at":"0001-01-01T00:00:00Z"}]`)

		printTablePodsList(body)
	})

	t.Run("error", func(t *testing.T) {
		body := []byte(`[{"namespace":"default"`)

		printTablePodsList(body)
	})
}

func TestCluster_print_table_pods_delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body := []byte(`{"pods":["container \"nginx-test-62bxfrkfgk-nginx-container\" in namespace \"scalezilla\": not found","container \"redis-test-3chnfcmbcn-redis-container\" in namespace \"scalezilla\": not found"]}`)

		printTablePodsDelete(body)
	})

	t.Run("error", func(t *testing.T) {
		body := []byte(`[{"namespace":"default"`)

		printTablePodsDelete(body)
	})
}

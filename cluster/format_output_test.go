package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_print_table(t *testing.T) {
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

package cluster

import (
	"testing"
)

func TestCluster_print_table(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		body := []byte(`[{"id":"127.0.0.1:20001","name":"127.0.0.1:20001","address":"127.0.0.1:20001","kind":"server","leader":true,"nodePool":"default"},{"id":"127.0.0.1:20005","name":"fake","address":"127.0.0.1:20005","kind":"client","leader":false,"nodePool":"default"},{"id":"127.0.0.1:20004","name":"127.0.0.1:20004","address":"127.0.0.1:20004","kind":"server","leader":false,"nodePool":"default"},{"id":"127.0.0.1:20007","name":"127.0.0.1:20007","address":"127.0.0.1:20007","kind":"server","leader":false,"nodePool":"default"}]`)

		// 		result := `ID               NAME             ADDRESS          KIND    LEADER  NODEPOOL
		// 127.0.0.1:20001  127.0.0.1:20001  127.0.0.1:20001  server  true    default`

		printTableNodesList(body)
	})

	t.Run("error", func(t *testing.T) {
		body := []byte(`[{"id":"127.0.0.1:20001"`)

		printTableNodesList(body)
	})
}

package cluster

import (
	"testing"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_receiver_responses(t *testing.T) {
	assert := assert.New(t)

	t.Run("rcv_service_ports_discovery", func(t *testing.T) {
		cfg := basicClusterConfig{randomPort: true, dev: true}
		cluster := makeBasicCluster(cfg)
		// cluster.buildAddressAndID()
		// defer func() {
		// 	_ = os.RemoveAll(cluster.config.DataDir)
		// }()

		responseChan := make(chan RPCResponse, 1)
		request := RPCRequest{
			RPCType: ServicePortsDiscovery,
			Request: &scalezillapb.ServicePortsDiscoveryRequestReply{
				Address:  "1234",
				Id:       "1234",
				PortHttp: uint32(defaultHTTPPort),
				PortGrpc: uint32(defaultGRPCPort),
				PortRaft: uint32(defaultRaftGRPCPort),
				IsVoter:  cluster.isVoter,
				NodePool: cluster.nodePool,
			},
			ResponseChan: responseChan,
		}

		go cluster.rcvServicePortsDiscovery(request)
		response := <-responseChan
		assert.NotNil(response.Response)
	})
}

package cluster

import (
	"testing"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"github.com/stretchr/testify/assert"
)

func TestCluster_grpc_content_builders(t *testing.T) {
	assert := assert.New(t)

	t.Run("service_ports_discovery", func(t *testing.T) {
		request := RPCServicePortsDiscoveryRequest{}
		response := RPCServicePortsDiscoveryResponse{}
		assert.Equal(&scalezillapb.ServicePortsDiscoveryRequestReply{}, makeServicePortsDiscoveryRequest(request))
		assert.Equal(response, makeServicePortsDiscoveryResponse(nil))

		request = RPCServicePortsDiscoveryRequest{Address: "X"}
		response = RPCServicePortsDiscoveryResponse{Address: "X"}
		assert.Equal(&scalezillapb.ServicePortsDiscoveryRequestReply{Address: "X"}, makeServicePortsDiscoveryRequest(request))
		assert.Equal(response, makeServicePortsDiscoveryResponse(&scalezillapb.ServicePortsDiscoveryRequestReply{Address: "X"}))
	})

	t.Run("service_node_polling", func(t *testing.T) {
		request := RPCServiceNodePollingRequest{}
		response := RPCServiceNodePollingResponse{}
		assert.Equal(&scalezillapb.ServiceNodePollingRequestReply{}, makeServiceNodePollingRequest(request))
		assert.Equal(response, makeServiceNodePollingResponse(nil))

		request = RPCServiceNodePollingRequest{Address: "X"}
		response = RPCServiceNodePollingResponse{Address: "X"}
		assert.Equal(&scalezillapb.ServiceNodePollingRequestReply{Address: "X"}, makeServiceNodePollingRequest(request))
		assert.Equal(response, makeServiceNodePollingResponse(&scalezillapb.ServiceNodePollingRequestReply{Address: "X"}))
	})
}

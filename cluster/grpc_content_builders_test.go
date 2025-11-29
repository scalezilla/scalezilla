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
}

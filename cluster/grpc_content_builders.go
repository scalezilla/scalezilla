package cluster

import "github.com/scalezilla/scalezilla/scalezillapb"

// makeServicePortsDiscoveryRequest build request for that matter
func makeServicePortsDiscoveryRequest(data RPCServicePortsDiscoveryRequest) *scalezillapb.ServicePortsDiscoveryRequestReply {
	return &scalezillapb.ServicePortsDiscoveryRequestReply{
		Address:  data.Address,
		Id:       data.ID,
		PortHttp: data.PortHTTP,
		PortGrpc: data.PortGRPC,
		PortRaft: data.PortRaft,
		IsVoter:  data.IsVoter,
		NodePool: data.NodePool,
	}
}

// makeServicePortsDiscoveryResponse build response for that matter
func makeServicePortsDiscoveryResponse(data *scalezillapb.ServicePortsDiscoveryRequestReply) RPCServicePortsDiscoveryResponse {
	if data == nil {
		return RPCServicePortsDiscoveryResponse{}
	}
	return RPCServicePortsDiscoveryResponse{
		Address:  data.Address,
		ID:       data.Id,
		PortHTTP: data.PortHttp,
		PortGRPC: data.PortGrpc,
		PortRaft: data.PortRaft,
		IsVoter:  data.IsVoter,
		NodePool: data.NodePool,
	}
}

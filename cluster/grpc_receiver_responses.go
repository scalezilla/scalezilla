package cluster

import "github.com/scalezilla/scalezilla/scalezillapb"

// rcvServicePortsDiscovery will answer back to rpc request from grpc receiver
func (c *Cluster) rcvServicePortsDiscovery(data RPCRequest) {
	request := data.Request.(*scalezillapb.ServicePortsDiscoveryRequestReply)
	c.mu.Lock()
	c.nodeMap[request.Id] = &nodeMap{
		Address:   request.Address,
		ID:        request.Id,
		HTTPPort:  request.PortHttp,
		GRPCPort:  request.PortGrpc,
		RaftyPort: request.PortRaft,
		IsVoter:   request.IsVoter,
		NodePool:  request.NodePool,
	}
	c.mu.Unlock()

	data.ResponseChan <- RPCResponse{
		Response: &scalezillapb.ServicePortsDiscoveryRequestReply{
			Address:  c.config.HostIPAddress,
			Id:       c.id,
			PortHttp: uint32(c.config.HTTPPort),
			PortGrpc: uint32(c.config.GRPCPort),
			PortRaft: uint32(c.config.RaftGRPCPort),
			IsVoter:  c.isVoter,
			NodePool: c.nodePool,
		},
	}
}

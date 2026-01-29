package cluster

// respServicePortsDiscovery will receive response from
// reqServicePortsDiscovery
func (c *Cluster) respServicePortsDiscovery(data RPCResponse) {
	if data.Error != nil {
		return
	}

	response := data.Response.(RPCServicePortsDiscoveryResponse)
	if response != (RPCServicePortsDiscoveryResponse{}) {
		c.logger.Debug().Msgf("discovery response address %s id %s node pool %s http port %d grpc port %d raft port %d\n", response.Address, response.ID, response.NodePool, response.PortHTTP, response.PortGRPC, response.PortRaft)
		c.mu.Lock()
		if _, ok := c.nodeMap[response.ID]; !ok {
			c.bootstrapExpectedSize.Add(1)
			c.nodeMap[response.ID] = &nodeMap{
				IsVoter:   response.IsVoter,
				ID:        response.ID,
				Address:   response.Address,
				HTTPPort:  response.PortHTTP,
				GRPCPort:  response.PortGRPC,
				RaftyPort: response.PortRaft,
				NodePool:  response.NodePool,
			}
		}
		c.mu.Unlock()
	}
}

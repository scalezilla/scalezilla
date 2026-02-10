package cluster

import "github.com/scalezilla/scalezilla/osdiscovery"

// respServicePortsDiscovery will receive response
// from reqServicePortsDiscovery
func (c *Cluster) respServicePortsDiscovery(data RPCResponse) {
	if data.Error != nil {
		return
	}

	response := data.Response.(RPCServicePortsDiscoveryResponse)
	if response != (RPCServicePortsDiscoveryResponse{}) {
		c.logger.Debug().Msgf("discovery response address %s id %s node pool %s http port %d grpc port %d raft port %d\n", response.Address, response.ID, response.NodePool, response.PortHTTP, response.PortGRPC, response.PortRaft)
		c.nodeMapMu.Lock()
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
		c.nodeMapMu.Unlock()
	}
}

// respServiceNodePolling will receive response
// from reqServicePNodePolling
func (c *Cluster) respServiceNodePolling(data RPCResponse) {
	if data.Error != nil {
		return
	}

	response, ok := data.Response.(RPCServiceNodePollingResponse)
	if ok {
		c.nodeMapMu.Lock()
		if _, ok := c.nodeMap[response.ID]; !ok {
			c.nodeMap[response.ID] = &nodeMap{
				Address: response.Address,
				ID:      response.ID,
			}
			c.nodeMap[response.ID].SystemInfo.OS = &osdiscovery.OS{
				Name:         response.OsName,
				Vendor:       response.OsVendor,
				Version:      response.OsVersion,
				Family:       response.OsFamily,
				Hostname:     response.OsHostname,
				Architecture: response.OsArchitecture,
				OSType:       response.OsType,
			}
			c.nodeMap[response.ID].SystemInfo.CPU = &osdiscovery.CPU{
				CPU:                 response.CpuTotal,
				Cores:               response.CpuCores,
				Frequency:           response.CpuFrequency,
				CumulativeFrequency: response.CpuCumulativeFrequency,
				Capabilitites:       response.CpuCapabilitites,
				Vendor:              response.CpuVendor,
				Model:               response.CpuModel,
			}
			c.nodeMap[response.ID].SystemInfo.Memory = &osdiscovery.Memory{
				Total:     uint(response.MemoryTotal),
				Available: uint(response.MemoryAvailable),
			}
			c.nodeMap[response.ID].Metadata = response.Metadata
		}
		c.nodeMapMu.Unlock()
	}
}

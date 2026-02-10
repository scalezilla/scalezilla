package cluster

import (
	"github.com/scalezilla/scalezilla/osdiscovery"
	"github.com/scalezilla/scalezilla/scalezillapb"
)

// rcvServicePortsDiscovery will answer back to rpc request from grpc receiver
func (c *Cluster) rcvServicePortsDiscovery(data RPCRequest) {
	request := data.Request.(*scalezillapb.ServicePortsDiscoveryRequestReply)
	c.nodeMapMu.Lock()
	if _, ok := c.nodeMap[request.Id]; !ok {
		if c.isVoter && !c.bootstrapExpectedSizeReach.Load() {
			c.bootstrapExpectedSize.Add(1)
		}
		if !c.isVoter {
			c.clientContactedServer.Store(true)
		}
		c.nodeMap[request.Id] = &nodeMap{
			Address:   request.Address,
			ID:        request.Id,
			HTTPPort:  request.PortHttp,
			GRPCPort:  request.PortGrpc,
			RaftyPort: request.PortRaft,
			IsVoter:   request.IsVoter,
			NodePool:  request.NodePool,
		}
	}
	c.nodeMapMu.Unlock()

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

// rcvServiceNodePolling will answer back to rpc request from grpc receiver
func (c *Cluster) rcvServiceNodePolling(data RPCRequest) {
	request := data.Request.(*scalezillapb.ServiceNodePollingRequestReply)
	c.nodeMapMu.Lock()
	if _, ok := c.nodeMap[request.Id]; !ok {
		c.nodeMap[request.Id] = &nodeMap{
			Address: request.Address,
			ID:      request.Id,
		}
		c.nodeMap[request.Id].SystemInfo.OS = &osdiscovery.OS{
			Name:         request.OsName,
			Vendor:       request.OsVendor,
			Version:      request.OsVersion,
			Family:       request.OsFamily,
			Hostname:     request.OsHostname,
			Architecture: request.OsArchitecture,
			OSType:       request.OsType,
		}
		c.nodeMap[request.Id].SystemInfo.CPU = &osdiscovery.CPU{
			CPU:                 request.CpuTotal,
			Cores:               request.CpuCores,
			Frequency:           request.CpuFrequency,
			CumulativeFrequency: request.CpuCumulativeFrequency,
			Capabilitites:       request.CpuCapabilitites,
			Vendor:              request.CpuVendor,
			Model:               request.CpuModel,
		}
		c.nodeMap[request.Id].SystemInfo.Memory = &osdiscovery.Memory{
			Total:     uint(request.MemoryTotal),
			Available: uint(request.MemoryAvailable),
		}
		c.nodeMap[request.Id].Metadata = request.Metadata
	}
	c.nodeMapMu.Unlock()

	response := &scalezillapb.ServiceNodePollingRequestReply{
		Address: c.config.HostIPAddress,
		Id:      c.id,
	}
	if c.systemInfo != nil {
		response.OsName = c.systemInfo.OS.Name
		response.OsVendor = c.systemInfo.OS.Vendor
		response.OsVersion = c.systemInfo.OS.Version
		response.OsFamily = c.systemInfo.OS.Family
		response.OsHostname = c.systemInfo.OS.Hostname
		response.OsArchitecture = c.systemInfo.OS.Architecture
		response.OsType = c.systemInfo.OS.OSType
		response.CpuTotal = c.systemInfo.CPU.CPU
		response.CpuCores = c.systemInfo.CPU.Cores
		response.CpuFrequency = c.systemInfo.CPU.Frequency
		response.CpuCumulativeFrequency = c.systemInfo.CPU.CumulativeFrequency
		response.CpuCapabilitites = c.systemInfo.CPU.Capabilitites
		response.CpuVendor = c.systemInfo.CPU.Vendor
		response.CpuModel = c.systemInfo.CPU.Model
		response.MemoryTotal = uint64(c.systemInfo.Memory.Total)
		response.MemoryAvailable = uint64(c.systemInfo.Memory.Available)
	}
	data.ResponseChan <- RPCResponse{
		Response: response,
	}
}

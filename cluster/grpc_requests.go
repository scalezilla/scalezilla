package cluster

import (
	"context"
	"time"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"google.golang.org/grpc"
)

// sendRPC is used to send rpc requests
func (c *Cluster) sendRPC(address string, client scalezillapb.ScalezillaClient, request RPCRequest) {
	options := []grpc.CallOption{}

	ctx := context.Background()
	var cancel context.CancelFunc
	if request.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), request.Timeout)
		defer cancel()
	}

	switch request.RPCType {
	case ServicePortsDiscovery:
		resp, err := client.ServicePortsDiscovery(
			ctx,
			makeServicePortsDiscoveryRequest(request.Request.(RPCServicePortsDiscoveryRequest)),
			options...,
		)
		request.ResponseChan <- RPCResponse{Response: makeServicePortsDiscoveryResponse(resp), Error: err, TargetNode: address}

	case ServiceNodePolling:
		resp, err := client.ServiceNodePolling(
			ctx,
			makeServiceNodePollingRequest(request.Request.(RPCServiceNodePollingRequest)),
			options...,
		)
		request.ResponseChan <- RPCResponse{Response: makeServiceNodePollingResponse(resp), Error: err, TargetNode: address}
	}
}

// reqServicePortsDiscovery will send a rpc request to other nodes
func (c *Cluster) reqServicePortsDiscovery() {
	request := RPCRequest{
		RPCType: ServicePortsDiscovery,
		Request: RPCServicePortsDiscoveryRequest{
			Address:  c.config.HostIPAddress,
			ID:       c.id,
			NodePool: c.nodePool,
			PortHTTP: uint32(c.config.HTTPPort),
			PortGRPC: uint32(c.config.GRPCPort),
			PortRaft: uint32(c.config.RaftGRPCPort),
			IsVoter:  c.isVoter,
		},
		Timeout:      time.Second,
		ResponseChan: c.rpcServicePortsDiscoveryChanResp,
	}

	for _, member := range c.members {
		go func() {
			if client := c.getClient(member); client != nil {
				c.di.sendRPCFunc(member, client, request)
			}
		}()
	}
}

// reqServiceNodePolling will send a rpc request to other nodes
func (c *Cluster) reqServiceNodePolling() {
	_ = c.checkSystemInfo()
	request := RPCRequest{
		RPCType: ServiceNodePolling,
		Request: RPCServiceNodePollingRequest{
			Address:                c.config.HostIPAddress,
			ID:                     c.id,
			OsName:                 c.systemInfo.OS.Name,
			OsVendor:               c.systemInfo.OS.Vendor,
			OsVersion:              c.systemInfo.OS.Version,
			OsFamily:               c.systemInfo.OS.Family,
			OsHostname:             c.systemInfo.OS.Hostname,
			OsArchitecture:         c.systemInfo.OS.Architecture,
			OsType:                 c.systemInfo.OS.OSType,
			CpuTotal:               c.systemInfo.CPU.CPU,
			CpuCores:               c.systemInfo.CPU.Cores,
			CpuFrequency:           c.systemInfo.CPU.Frequency,
			CpuCumulativeFrequency: c.systemInfo.CPU.CumulativeFrequency,
			CpuCapabilitites:       c.systemInfo.CPU.Capabilitites,
			CpuVendor:              c.systemInfo.CPU.Vendor,
			CpuModel:               c.systemInfo.CPU.Model,
			MemoryTotal:            uint64(c.systemInfo.Memory.Total),
			MemoryAvailable:        uint64(c.systemInfo.Memory.Available),
			Metadata:               c.config.Metadata,
		},
		Timeout:      time.Second,
		ResponseChan: c.rpcServiceNodePollingChanResp,
	}

	for _, member := range c.members {
		go func() {
			if client := c.getClient(member); client != nil {
				c.di.sendRPCFunc(member, client, request)
			}
		}()
	}
}

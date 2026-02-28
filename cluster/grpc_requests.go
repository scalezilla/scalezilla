package cluster

import (
	"context"
	"fmt"
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

	case ServiceNodeRegister:
		resp, err := client.ServiceNodeRegister(
			ctx,
			makeServiceNodeRegisterRequest(request.Request.(RPCServiceNodeRegisterRequest)),
			options...,
		)
		request.ResponseChan <- RPCResponse{Response: makeServiceNodeRegisterResponse(resp), Error: err, TargetNode: address}
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
			Members:  c.members_grpc,
		},
		Timeout:      time.Second,
		ResponseChan: c.rpcServicePortsDiscoveryChanResp,
	}

	for _, member := range c.members_grpc {
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
			OsHostname:             c.config.Hostname,
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

	for _, member := range c.members_grpc {
		go func() {
			if client := c.getClient(member); client != nil {
				c.di.sendRPCFunc(member, client, request)
			}
		}()
	}
}

// reqServiceNodeRegister will send a rpc request to the leader
// to be part of the cluster
func (c *Cluster) reqServiceNodeRegister() {
	request := RPCRequest{
		RPCType: ServiceNodeRegister,
		Request: RPCServiceNodeRegisterRequest{
			Address: c.raftyAddress.String(),
			ID:      c.id,
			IsVoter: c.isVoter,
		},
		Timeout:      time.Second,
		ResponseChan: c.rpcServiceNodeRegisterChanResp,
	}

	if ok, leaderAddress, leaderId := c.rafty.FetchLeader(); ok {
		if grpcAddress, ok := c.getServiceAddressFromRaft(leaderAddress); ok {
			c.logger.Debug().
				Str("address", c.raftyAddress.String()).
				Str("id", c.id).
				Str("isVoter", fmt.Sprintf("%t", c.isVoter)).
				Str("leaderAddress", grpcAddress).
				Str("leaderId", leaderId).
				Msgf("Asking for cluster membership")
			if client := c.getClient(grpcAddress); client != nil {
				c.di.sendRPCFunc(grpcAddress, client, request)
			}
		}
	}
}

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
	}
}

// reqServicePortsDiscovery will send a rpc request to other nodes
func (c *Cluster) reqServicePortsDiscovery() {
	request := RPCRequest{
		RPCType: ServicePortsDiscovery,
		Request: RPCServicePortsDiscoveryRequest{
			Address:  c.grpcAddress.String(),
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

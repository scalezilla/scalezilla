package cluster

import (
	"context"
	"time"

	"github.com/scalezilla/scalezilla/scalezillapb"
)

// ServicePortsDiscovery is used to receive rpc calls from
// other nodes
func (c *Cluster) ServicePortsDiscovery(ctx context.Context, in *scalezillapb.ServicePortsDiscoveryRequestReply) (*scalezillapb.ServicePortsDiscoveryRequestReply, error) {
	responseChan := make(chan RPCResponse, 1)
	select {
	case c.rpcServicePortsDiscoveryChanReq <- RPCRequest{
		RPCType:      ServicePortsDiscovery,
		Request:      in,
		ResponseChan: responseChan,
	}:

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-c.ctx.Done():
		return nil, ErrShutdown

	case <-time.After(500 * time.Millisecond):
		return nil, ErrTimeout
	}

	select {
	case response := <-responseChan:
		return response.Response.(*scalezillapb.ServicePortsDiscoveryRequestReply), response.Error

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-c.ctx.Done():
		return nil, ErrShutdown

	case <-time.After(time.Second):
		return nil, ErrTimeout
	}
}

// ServiceNodePolling is used to receive rpc calls
// from other nodes
func (c *Cluster) ServiceNodePolling(ctx context.Context, in *scalezillapb.ServiceNodePollingRequestReply) (*scalezillapb.ServiceNodePollingRequestReply, error) {
	responseChan := make(chan RPCResponse, 1)
	select {
	case c.rpcServiceNodePollingChanReq <- RPCRequest{
		RPCType:      ServiceNodePolling,
		Request:      in,
		ResponseChan: responseChan,
	}:

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-c.ctx.Done():
		return nil, ErrShutdown

	case <-time.After(500 * time.Millisecond):
		return nil, ErrTimeout
	}

	select {
	case response := <-responseChan:
		return response.Response.(*scalezillapb.ServiceNodePollingRequestReply), response.Error

	case <-ctx.Done():
		return nil, ctx.Err()

	case <-c.ctx.Done():
		return nil, ErrShutdown

	case <-time.After(time.Second):
		return nil, ErrTimeout
	}
}

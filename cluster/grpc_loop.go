package cluster

import "time"

// grpcLoop receive all rpc requests or responses
// from other nodes and also client commands
func (c *Cluster) grpcLoop() {
	for {
		select {
		case <-c.ctx.Done():
			c.drainRCVServicePortsDiscovery()
			c.drainRespServicePortsDiscovery()
			c.drainRCVServiceNodePolling()
			c.drainRespServiceNodePolling()
			return

		case data, ok := <-c.rpcServicePortsDiscoveryChanReq:
			if ok {
				c.rcvServicePortsDiscovery(data)
			}

		case data, ok := <-c.rpcServicePortsDiscoveryChanResp:
			if ok {
				c.respServicePortsDiscovery(data)
			}

		case <-time.After(c.nodePollingTimer):
			if c.isRunning.Load() && !c.dev {
				c.reqServiceNodePolling()
			}

		case data, ok := <-c.rpcServiceNodePollingChanReq:
			if ok {
				c.rcvServiceNodePolling(data)
			}

		case data, ok := <-c.rpcServiceNodePollingChanResp:
			if ok {
				c.respServiceNodePolling(data)
			}
		}
	}
}

// drainRCVServicePortsDiscovery will drain all remaining data in the chan
func (c *Cluster) drainRCVServicePortsDiscovery() {
	for {
		select {
		case data := <-c.rpcServicePortsDiscoveryChanReq:
			select {
			case data.ResponseChan <- RPCResponse{
				Error: ErrShutdown,
			}:
			//nolint staticcheck
			default:
			}
		//nolint staticcheck
		default:
			return
		}
	}
}

// drainRespServicePortsDiscovery will drain all remaining data in the chan
func (c *Cluster) drainRespServicePortsDiscovery() {
	for {
		select {
		case <-c.rpcServicePortsDiscoveryChanResp:
		//nolint staticcheck
		default:
			return
		}
	}
}

// rpcServiceNodePollingChanReq will drain all remaining data in the chan
func (c *Cluster) drainRCVServiceNodePolling() {
	for {
		select {
		case data := <-c.rpcServiceNodePollingChanReq:
			select {
			case data.ResponseChan <- RPCResponse{
				Error: ErrShutdown,
			}:
			//nolint staticcheck
			default:
			}
		//nolint staticcheck
		default:
			return
		}
	}
}

// drainRespServiceNodePolling will drain all remaining data in the chan
func (c *Cluster) drainRespServiceNodePolling() {
	for {
		select {
		case <-c.rpcServiceNodePollingChanResp:
		//nolint staticcheck
		default:
			return
		}
	}
}

package cluster

// grpcLoop receive all rpc requests or responses
// from other nodes and also client commands
func (c *Cluster) grpcLoop() {
	for {
		select {
		case <-c.ctx.Done():
			c.drainRCVServicePortsDiscovery()
			c.drainRespServicePortsDiscovery()
			return

		case data, ok := <-c.rpcServicePortsDiscoveryChanReq:
			if ok {
				c.rcvServicePortsDiscovery(data)
			}

		case data, ok := <-c.rpcServicePortsDiscoveryChanResp:
			if ok {
				c.respServicePortsDiscovery(data)
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

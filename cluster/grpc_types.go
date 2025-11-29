package cluster

import "time"

// RPCType is used to build rpc requests
type RPCType uint8

const (
	// ServicePortsDiscovery will be used to perform service ports discovery
	ServicePortsDiscovery RPCType = iota
)

// RPCRequest is used by chans in order to manage rpc requests
type RPCRequest struct {
	RPCType      RPCType
	Request      any
	Timeout      time.Duration
	ResponseChan chan<- RPCResponse
}

// RPCResponse  is used by RPCRequest in order to reply to rpc requests
type RPCResponse struct {
	TargetNode string
	Response   any
	Error      error
}

// RPCServicePortsDiscoveryRequest holds the requirements to perform service ports discovery
type RPCServicePortsDiscoveryRequest struct {
	Address, ID, NodePool        string
	PortHTTP, PortGRPC, PortRaft uint32
	IsVoter                      bool
}

// RPCServicePortsDiscoveryResponse holds the response from RPCServicePortsDiscoveryRequest
type RPCServicePortsDiscoveryResponse struct {
	Address, ID, NodePool        string
	PortHTTP, PortGRPC, PortRaft uint32
	IsVoter                      bool
}

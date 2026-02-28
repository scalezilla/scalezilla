package cluster

import (
	"time"
)

// RPCType is used to build rpc requests
type RPCType uint8

const (
	// ServicePortsDiscovery will be used to perform service ports discovery
	ServicePortsDiscovery RPCType = iota

	// ServiceNodePolling will be used to perform service node polling
	ServiceNodePolling

	// ServiceNodeRegister will be used to register node
	ServiceNodeRegister
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
	Members                      []string
}

// RPCServicePortsDiscoveryResponse holds the response from RPCServicePortsDiscoveryRequest
type RPCServicePortsDiscoveryResponse struct {
	Address, ID, NodePool        string
	PortHTTP, PortGRPC, PortRaft uint32
	IsVoter                      bool
	Members                      []string
}

// RPCServiceNodePollingRequest holds the requirements to perform node polling
type RPCServiceNodePollingRequest struct {
	Address, ID                           string
	OsName, OsVendor, OsVersion, OsFamily string
	OsHostname, OsArchitecture, OsType    string
	CpuTotal, CpuCores                    uint32
	CpuFrequency                          float32
	CpuCumulativeFrequency                float64
	CpuCapabilitites                      []string
	CpuVendor, CpuModel                   string
	MemoryTotal, MemoryAvailable          uint64
	Metadata                              map[string]string
}

// RPCServiceNodePollingResponse holds the response from RPCServiceNodePollingRequest
type RPCServiceNodePollingResponse struct {
	Address, ID                           string
	OsName, OsVendor, OsVersion, OsFamily string
	OsHostname, OsArchitecture, OsType    string
	CpuTotal, CpuCores                    uint32
	CpuFrequency                          float32
	CpuCumulativeFrequency                float64
	CpuCapabilitites                      []string
	CpuVendor, CpuModel                   string
	MemoryTotal, MemoryAvailable          uint64
	Metadata                              map[string]string
}

// RPCServiceNodeRegisterRequest holds the requirements
// to be part of the cluster
type RPCServiceNodeRegisterRequest struct {
	Address, ID string
	IsVoter     bool
}

// RPCServiceNodeRegisterResponse holds the response from RPCServiceNodeRegisterRequest
type RPCServiceNodeRegisterResponse struct {
	Acknowledged bool
}

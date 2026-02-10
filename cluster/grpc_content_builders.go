package cluster

import (
	"github.com/scalezilla/scalezilla/scalezillapb"
)

// makeServicePortsDiscoveryRequest build request for that matter
func makeServicePortsDiscoveryRequest(data RPCServicePortsDiscoveryRequest) *scalezillapb.ServicePortsDiscoveryRequestReply {
	return &scalezillapb.ServicePortsDiscoveryRequestReply{
		Address:  data.Address,
		Id:       data.ID,
		PortHttp: data.PortHTTP,
		PortGrpc: data.PortGRPC,
		PortRaft: data.PortRaft,
		IsVoter:  data.IsVoter,
		NodePool: data.NodePool,
	}
}

// makeServicePortsDiscoveryResponse build response for that matter
func makeServicePortsDiscoveryResponse(data *scalezillapb.ServicePortsDiscoveryRequestReply) RPCServicePortsDiscoveryResponse {
	if data == nil {
		return RPCServicePortsDiscoveryResponse{}
	}
	return RPCServicePortsDiscoveryResponse{
		Address:  data.Address,
		ID:       data.Id,
		PortHTTP: data.PortHttp,
		PortGRPC: data.PortGrpc,
		PortRaft: data.PortRaft,
		IsVoter:  data.IsVoter,
		NodePool: data.NodePool,
	}
}

// makeServiceNodePollingRequest build request for that matter
func makeServiceNodePollingRequest(data RPCServiceNodePollingRequest) *scalezillapb.ServiceNodePollingRequestReply {
	return &scalezillapb.ServiceNodePollingRequestReply{
		Address:                data.Address,
		Id:                     data.ID,
		OsName:                 data.OsName,
		OsVendor:               data.OsVendor,
		OsVersion:              data.OsVersion,
		OsFamily:               data.OsFamily,
		OsHostname:             data.OsHostname,
		OsArchitecture:         data.OsArchitecture,
		OsType:                 data.OsType,
		CpuTotal:               data.CpuTotal,
		CpuCores:               data.CpuCores,
		CpuFrequency:           data.CpuFrequency,
		CpuCumulativeFrequency: data.CpuCumulativeFrequency,
		CpuCapabilitites:       data.CpuCapabilitites,
		CpuVendor:              data.CpuVendor,
		CpuModel:               data.CpuModel,
		MemoryTotal:            uint64(data.MemoryTotal),
		MemoryAvailable:        uint64(data.MemoryAvailable),
	}
}

// makeServicePortsDiscoveryResponse build response for that matter
func makeServiceNodePollingResponse(data *scalezillapb.ServiceNodePollingRequestReply) RPCServiceNodePollingResponse {
	if data == nil {
		return RPCServiceNodePollingResponse{}
	}
	return RPCServiceNodePollingResponse{
		Address:                data.Address,
		ID:                     data.Id,
		OsName:                 data.OsName,
		OsVendor:               data.OsVendor,
		OsVersion:              data.OsVersion,
		OsFamily:               data.OsFamily,
		OsHostname:             data.OsHostname,
		OsArchitecture:         data.OsArchitecture,
		OsType:                 data.OsType,
		CpuTotal:               data.CpuTotal,
		CpuCores:               data.CpuCores,
		CpuFrequency:           data.CpuFrequency,
		CpuCumulativeFrequency: data.CpuCumulativeFrequency,
		CpuCapabilitites:       data.CpuCapabilitites,
		CpuVendor:              data.CpuVendor,
		CpuModel:               data.CpuModel,
		MemoryTotal:            uint64(data.MemoryTotal),
		MemoryAvailable:        uint64(data.MemoryAvailable),
	}
}

package cluster

import (
	"fmt"
	"maps"
	"net"
	"slices"
	"strings"
)

// getServiceAddressFromRaft returns the grpc service address
// from raft address if found from node map
func (c *Cluster) getServiceAddressFromRaft(address string) (string, bool) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return "", false
	}

	c.nodeMapMu.RLock()
	keys := slices.Sorted(maps.Keys(c.nodeMap))
	if index := slices.IndexFunc(keys, func(p string) bool {
		return strings.Contains(p, host)
	}); index != -1 {
		if m, ok := c.nodeMap[keys[index]]; ok {
			return fmt.Sprintf("%s:%d", host, m.GRPCPort), true
		}
	}
	c.nodeMapMu.RUnlock()
	return "", false
}

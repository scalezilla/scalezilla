package cluster

import (
	"net/http"
	"sort"
	"strings"

	"github.com/Lord-Y/rafty"
	"github.com/gin-gonic/gin"
)

// nodesList is used to list cluster nodes
// with user provided token
func (cc *Cluster) nodesList(c *gin.Context) {
	var z []APINodesListResponse
	if cc.dev {
		z = append(z, APINodesListResponse{
			ID:       cc.id,
			Name:     cc.systemInfo.OS.Hostname,
			Address:  cc.raftyAddress.String(),
			Leader:   true,
			Kind:     "server",
			NodePool: cc.nodePool,
		})
		c.JSON(http.StatusOK, z)
		return
	}

	if !cc.rafty.IsBootstrapped() {
		c.JSON(http.StatusForbidden, gin.H{"error": ErrClusterNotBootstrapped.Error()})
		return
	}

	var req APINodesListRequest
	_ = c.Bind(&req)

	hasLeader, _, leaderId := cc.rafty.Leader()
	if hasLeader {
		status := cc.rafty.Status()
		cc.nodeMapMu.RLock()
		for _, v := range status.Configuration.ServerMembers {
			voter := "client"
			if v.IsVoter {
				voter = "server"
			}
			var (
				hostname, nodePool string
				leader             bool
			)
			if m, ok := cc.nodeMap[v.ID]; ok {
				if m.SystemInfo.OS != nil {
					hostname = m.SystemInfo.OS.Hostname
				}
				nodePool = m.NodePool
			}
			if v.ID == cc.id {
				nodePool = cc.nodePool
				hostname = cc.config.Hostname
			}

			// the following must be stay like so to make sure we have
			// the right leader on the right node
			if v.ID == leaderId {
				leader = true
			}

			switch req.Kind {
			case voter:
				z = append(z, APINodesListResponse{
					ID:       v.ID,
					Name:     hostname,
					Address:  v.Address,
					Leader:   leader,
					Kind:     voter,
					NodePool: nodePool,
				})
			case "":
				z = append(z, APINodesListResponse{
					ID:       v.ID,
					Name:     hostname,
					Address:  v.Address,
					Leader:   leader,
					Kind:     voter,
					NodePool: nodePool,
				})
			}
		}
		cc.nodeMapMu.RUnlock()
		if len(z) > 1 && (req.Kind == "" || req.Kind == "server") {
			sortNodes(z)
		}
		c.JSON(http.StatusOK, z)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": rafty.ErrNoLeader.Error()})
}

// sortNodes allow us to sort nodes by leader and kind
func sortNodes(nodes []APINodesListResponse) {
	sort.SliceStable(nodes, func(i, j int) bool {
		a, b := nodes[i], nodes[j]

		if a.Leader != b.Leader {
			return a.Leader && !b.Leader
		}

		aServer := strings.EqualFold(a.Kind, "server")
		bServer := strings.EqualFold(b.Kind, "server")
		if aServer != bServer {
			return aServer && !bServer
		}

		return a.ID < b.ID
	})
}

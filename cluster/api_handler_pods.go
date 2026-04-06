package cluster

import (
	"net/http"

	"github.com/Lord-Y/rafty"
	"github.com/gin-gonic/gin"
)

// deploymentApply is used to create deployment with user provided token
func (cc *Cluster) podsList(c *gin.Context) {
	var req APIPodsListRequest
	_ = c.ShouldBind(&req)

	if cc.rafty.IsLeader() {
		list, err := cc.di.listContainerFunc(cc.ctx, req.Namespace)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		var result []APIPodsListResponse
		if len(list) == 0 {
			c.JSON(http.StatusNotFound, result)
			return
		}

		for _, l := range list {
			result = append(result, APIPodsListResponse{
				Namespace: l.Namespace,
				ID:        l.ID,
				Image:     l.Image,
				PID:       l.PID,
				Runtime:   l.Runtime,
				Status:    l.Status,
				CreatedAt: l.CreatedAt,
			})
		}
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": rafty.ErrNotLeader.Error()})
}

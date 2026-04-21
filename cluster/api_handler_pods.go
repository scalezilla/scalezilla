package cluster

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/Lord-Y/rafty"
	"github.com/gin-gonic/gin"
)

// podsList is used to list pods with user provided informations
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

// podsDelete is used to list pods with user provided informations
func (cc *Cluster) podsDelete(c *gin.Context) {
	var req APIPodsDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Namespace == "" {
		req.Namespace = defaultPodNamespace
	}

	if cc.rafty.IsLeader() {
		var result APIPodsDeleteResponse
		if req.Detached {
			for _, pod := range req.Pods {
				result.Pods = append(result.Pods, fmt.Sprintf("pod %s deletion initiated in background", pod))
			}
			c.JSON(http.StatusOK, result)

			var wg sync.WaitGroup
			for _, pod := range req.Pods {
				wg.Go(func() {
					if err := cc.di.deleteContainerFunc(cc.ctx, req.Namespace, pod, defaultDeletePodTimeout); err != nil {
						cc.logger.Error().Err(err).
							Str("component", "pods").
							Str("namespace", req.Namespace).
							Str("pod", pod).
							Msg("Fail to delete pod")
					}
				})
			}
			wg.Wait()
			return
		}

		for _, pod := range req.Pods {
			if err := cc.di.deleteContainerFunc(cc.ctx, req.Namespace, pod, defaultDeletePodTimeout); err != nil {
				cc.logger.Error().Err(err).
					Str("component", "pods").
					Str("namespace", req.Namespace).
					Str("pod", pod).
					Msg("Fail to delete pod")

				result.Pods = append(result.Pods, strings.ReplaceAll(err.Error(), "scalezilla", req.Namespace))
			} else {
				result.Pods = append(result.Pods, fmt.Sprintf("pod %s deleted successfully", pod))
			}
		}

		// TODO: For later maybe if there is an error we should be return a different status code
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": rafty.ErrNotLeader.Error()})
}

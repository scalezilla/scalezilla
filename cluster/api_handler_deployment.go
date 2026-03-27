package cluster

import (
	"net/http"

	"github.com/Lord-Y/rafty"
	"github.com/gin-gonic/gin"
	"github.com/scalezilla/scalezilla/cri"
)

// deploymentApply is used to create deployment with user provided token
func (cc *Cluster) deploymentApply(c *gin.Context) {
	var req APIDeploymentApplyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrDeploymentPayload.Error()})
		return
	}

	if cc.rafty.IsLeader() {
		spec, err := cc.parseDeployment([]byte(req.HCLContent))
		if err != nil {
			cc.logger.Error().Err(err).Str("component", "deployment").Msgf("parsing error")
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		// this transformation is only temporary
		// as we need to mature the deployment
		cspec := cri.CreateContainerSpec{
			Namespace:   spec.Deployment.Namespace,
			ContainerID: spec.Deployment.Pod.Container.Name,
			Image: cri.ImageSpec{
				Image: spec.Deployment.Pod.Container.Image,
			},
			DefaultLogPath: "/var/log/containerd",
		}

		if err := cc.di.createContainerFunc(cc.ctx, cspec); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, APIGenericResponse{Message: "OK"})
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": rafty.ErrNotLeader.Error()})
}

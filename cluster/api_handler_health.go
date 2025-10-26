package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// health return the liveness of the cluster
func (cc *Cluster) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// healthz return the readiness of the cluster
func (cc *Cluster) healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

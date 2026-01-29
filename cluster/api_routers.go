package cluster

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// newApiRouters will return the api router
func (c *Cluster) newApiRouters() *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(requestid.New())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/cluster/health", c.health)
		v1.GET("/cluster/healthz", c.healthz)

		v1.GET("/cluster/bootstrap/status", c.bootstapStatus)
		v1.POST("/cluster/bootstrap/cluster", c.bootstrapCluster)
	}
	return router
}

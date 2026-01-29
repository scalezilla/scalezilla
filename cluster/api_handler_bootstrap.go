package cluster

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// bootstapStatus return if the cluster has been bootstrapped
func (cc *Cluster) bootstapStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"bootstrapped": cc.rafty.IsBootstrapped()})
}

// bootstrapCluster is used to bootstrap cluster
// with or without user provided token
func (cc *Cluster) bootstrapCluster(c *gin.Context) {
	if cc.rafty.IsBootstrapped() {
		c.JSON(http.StatusForbidden, gin.H{"error": "cluster already boostrapped"})
		return
	}

	var req APIBootstrapClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil && c.Request.Body != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tok string
	if req.Token == "" {
		tok = uuid.NewString()
	} else {
		tok = req.Token
	}
	z := &AclToken{
		AccessorID:   uuid.NewString(),
		Token:        tok,
		InitialToken: true,
	}
	if err := cc.submitCommandACLTokenWrite(aclTokenCommandSet, z); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, z)
}

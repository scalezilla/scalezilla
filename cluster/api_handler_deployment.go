package cluster

import (
	"crypto/rand"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

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

		var (
			replicaSetID string
			state        deploymentState
			version      uint64
		)
		if cc.fsm.memoryDeploymentExistsFunc([]byte(spec.Deployment.Namespace), []byte(spec.Deployment.Name)) {
			data, err := cc.fsm.memoryDeploymentGetFunc([]byte(spec.Deployment.Namespace), []byte(spec.Deployment.Name))
			if err != nil && err != rafty.ErrKeyNotFound {
				cc.logger.Error().Err(err).
					Str("component", "deployment").
					Str("namespace", spec.Deployment.Namespace).
					Str("deploymentName", spec.Deployment.Name).
					Msgf("fail to get deployment")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			state, err = deploymentDecodeCommand(data)
			if err != nil {
				cc.logger.Error().Err(err).
					Str("component", "deployment").
					Str("namespace", spec.Deployment.Namespace).
					Str("deploymentName", spec.Deployment.Name).
					Msgf("fail to decode retrieved deployment")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			if state.Content[state.CurrentUsedVersion].RawContent == req.HCLContent {
				c.JSON(http.StatusOK, APIGenericResponse{Message: "already deployed"})
				return
			} else {
				keys := slices.Sorted(maps.Keys(state.Content))
				version = keys[0] + 1
				replicaSetID = string(strings.ToLower(rand.Text())[:10])
				dc := deploymentContent{
					RawContent:   req.HCLContent,
					Version:      version,
					CreatedAt:    time.Now(),
					ReplicaSetID: replicaSetID,
				}
				state := deploymentState{
					Kind:              deploymentCommandSet,
					Name:              spec.Deployment.Name,
					NewRollingVersion: int64(version),
					MustBeStarted:     true,
				}
				state.Content = make(map[uint64]deploymentContent)
				state.Content[version] = dc
				if err := cc.di.submitCommandDeploymentWriteFunc(10*time.Second, state); err != nil {
					cc.logger.Error().Err(err).
						Str("component", "deployment").
						Str("namespace", spec.Deployment.Namespace).
						Str("deploymentName", spec.Deployment.Name).
						Msgf("fail to submit deployment")
					c.JSON(http.StatusInternalServerError, gin.H{"error": err})
					return
				}
			}
		} else {
			replicaSetID = string(strings.ToLower(rand.Text())[:10])
			version = 1
			dc := deploymentContent{
				RawContent:   req.HCLContent,
				Version:      version,
				CreatedAt:    time.Now(),
				ReplicaSetID: replicaSetID,
			}
			content := make(map[uint64]deploymentContent, 1)
			content[dc.Version] = dc
			state = deploymentState{
				Kind:               deploymentCommandSet,
				Name:               spec.Deployment.Name,
				NewRollingVersion:  int64(version),
				CurrentUsedVersion: version,
				Content:            content,
				MustBeStarted:      true,
			}
			if err := cc.di.submitCommandDeploymentWriteFunc(10*time.Second, state); err != nil {
				cc.logger.Error().Err(err).
					Str("component", "deployment").
					Str("namespace", spec.Deployment.Namespace).
					Str("deploymentName", spec.Deployment.Name).
					Msgf("fail to submit deployment")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
		}

		// this transformation is only temporary
		// as we need to mature the deployment
		cspec := cri.CreateContainerSpec{
			Namespace:   spec.Deployment.Namespace,
			ContainerID: fmt.Sprintf("%s-%s-%s", spec.Deployment.Name, replicaSetID, spec.Deployment.Pod.Container.Name),
			Image: cri.ImageSpec{
				Image: spec.Deployment.Pod.Container.Image,
			},
			DefaultLogPath: "/var/log/containerd",
			Labels:         spec.Deployment.Metadata,
		}

		if err := cc.di.createContainerFunc(cc.ctx, cspec); err != nil {
			cc.logger.Error().Err(err).
				Str("component", "deployment").
				Str("namespace", spec.Deployment.Namespace).
				Str("deploymentName", spec.Deployment.Name).
				Msgf("fail to create container")
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, APIGenericResponse{Message: "deployment successful"})
		// make version stable
		state.NewRollingVersion = -1
		if entry, ok := state.Content[version]; ok {
			entry.IsStable = true
			state.Content[version] = entry
			if err := cc.di.stableSubmitCommandDeploymentWriteFunc(10*time.Second, state); err != nil {
				cc.logger.Error().Err(err).
					Str("component", "deployment").
					Str("namespace", spec.Deployment.Namespace).
					Str("deploymentName", spec.Deployment.Name).
					Msgf("fail to set deployment as stable")
			}
		}
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": rafty.ErrNotLeader.Error()})
}

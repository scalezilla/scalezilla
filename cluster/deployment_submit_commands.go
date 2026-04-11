package cluster

import (
	"bytes"
	"time"

	"github.com/Lord-Y/rafty"
)

// submitCommandDeploymentWrite will send a write command to the leader
func (c *Cluster) submitCommandDeploymentWrite(timeout time.Duration, cmd deploymentState) error {
	buffer := new(bytes.Buffer)

	if err := c.di.deploymentEncodeCommandFunc(cmd, buffer); err != nil {
		return err
	}

	if _, err := c.rafty.SubmitCommand(timeout, rafty.LogReplication, buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

package cluster

import (
	"bytes"
	"time"

	"github.com/Lord-Y/rafty"
)

// submitCommandACLTokenWrite will send a write command to the leader
func (c *Cluster) submitCommandACLTokenWrite(kind commandKind, data *AclToken) error {
	buffer := new(bytes.Buffer)
	cmd := aclTokenCommand{
		Kind:         kind,
		AccessorID:   data.AccessorID,
		Token:        data.Token,
		InitialToken: data.InitialToken,
	}

	if err := c.di.aclTokenEncodeCommandFunc(cmd, buffer); err != nil {
		return err
	}

	if _, err := c.rafty.SubmitCommand(time.Second, rafty.LogReplication, buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

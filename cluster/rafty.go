package cluster

import (
	"time"

	"github.com/Lord-Y/rafty"
)

// newRafty is use to build rafty config
func (c *Cluster) newRafty() (*rafty.Rafty, error) {
	options := rafty.Options{
		Logger:                 c.logger,
		DataDir:                c.config.DataDir,
		ElectionTimeout:        5000,
		HeartbeatTimeout:       3000,
		IsVoter:                c.isVoter,
		InitialPeers:           c.buildPeers(),
		MetricsNamespacePrefix: c.raftMetricPrefix,
		ShutdownOnRemove:       true,
		BootstrapCluster:       true,
		MaxAppendEntries:       defaultMaxAppendEntries,
	}

	if c.isVoter {
		options.SnapshotInterval = c.config.Server.Raft.SnapshotInterval
		options.SnapshotThreshold = c.config.Server.Raft.SnapshotThreshold
		options.MinimumClusterSize = c.config.Server.Raft.BootstrapExpectedSize
	} else {
		options.SnapshotInterval = c.config.Client.Raft.SnapshotInterval
		options.SnapshotThreshold = c.config.Client.Raft.SnapshotThreshold
		options.LeaveOnTerminate = true
	}

	if c.dev {
		options.IsSingleServerCluster = c.dev
		options.BootstrapCluster = false
	}

	var err error
	if c.raftyStore, err = c.buildStore(); err != nil {
		return nil, err
	}

	c.fsm = newFSM(c.raftyStore)
	snapshotConfig := newSnapshot(c.config.DataDir, 3)

	return rafty.NewRafty(c.address, c.id, options, c.raftyStore, c.raftyStore, c.fsm, snapshotConfig)
}

// startRafty is use to start rafty cluster
func (c *Cluster) startRafty() error {
	errChan := make(chan error, 1)
	defer close(errChan)
	go func() {
		if err := c.rafty.Start(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(10 * time.Millisecond):
		return nil
	}
}

// stopRafty is used to stop rafty cluster
func (c *Cluster) stopRafty() {
	c.rafty.Stop()
}

// raftyStoreClose is used to close the store used by rafty
func (c *Cluster) raftyStoreClose() error {
	return c.raftyStore.Close()
}

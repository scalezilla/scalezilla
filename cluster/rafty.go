package cluster

import (
	"time"

	"github.com/Lord-Y/rafty"
)

// newRafty is use to build rafty config
func (c *Cluster) newRafty(metricsNamespacePrefix string) (*rafty.Rafty, error) {
	options := rafty.Options{
		Logger:                 c.logger,
		DataDir:                c.dataDir,
		MinimumClusterSize:     3,
		ElectionTimeout:        3000,
		HeartbeatTimeout:       1000,
		IsVoter:                true,
		InitialPeers:           c.buildPeers(),
		SnapshotInterval:       30 * time.Second,
		IsSingleServerCluster:  c.dev,
		MetricsNamespacePrefix: metricsNamespacePrefix,
	}

	var err error
	if c.raftyStore, err = c.buildStore(); err != nil {
		return nil, err
	}

	c.fsm = newFSM(c.raftyStore)
	snapshotConfig := newSnapshot(c.dataDir, 3)

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

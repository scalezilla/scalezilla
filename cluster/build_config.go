package cluster

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Lord-Y/rafty"
	bolt "go.etcd.io/bbolt"
)

// buildAddressAndID will build the node address and id
func (c *Cluster) buildAddressAndID() {
	c.grpcAddress = net.TCPAddr{
		IP:   net.ParseIP(c.config.HostIPAddress),
		Port: int(c.config.GRPCPort),
	}

	c.raftyAddress = net.TCPAddr{
		IP:   net.ParseIP(c.config.HostIPAddress),
		Port: int(c.config.RaftGRPCPort),
	}

	c.id = fmt.Sprintf("%s:%d", c.config.HostIPAddress, c.config.RaftGRPCPort)
}

// buildPeers will build the initial peer members of the cluster
func (c *Cluster) buildPeers() []rafty.InitialPeer {
	peers := []rafty.InitialPeer{{Address: c.raftyAddress.String()}}

	c.nodeMapMu.RLock()
	for _, v := range c.nodeMap {
		peers = append(peers, rafty.InitialPeer{Address: v.Address})
	}
	c.nodeMapMu.RUnlock()

	return peers
}

// buildDataDir will build the working dir of the current node
func (c *Cluster) buildDataDir() {
	c.config.DataDir = filepath.Join(os.TempDir(), "scalezilla", c.id)
}

// buildStore will build the bolt store
func (c *Cluster) buildStore() (*rafty.BoltStore, error) {
	storeOptions := rafty.BoltOptions{
		DataDir: c.config.DataDir,
		Options: bolt.DefaultOptions,
	}

	return rafty.NewBoltStorage(storeOptions)
}

// BuildSignal will build signal requirements based on provided
// context
func BuildSignal(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
}

// buildDevConfig build dev config to start a dev single cluster
func (c *Cluster) buildDevConfig(config ClusterInitialConfig) {
	c.isVoter = true
	c.dev = true
	c.raftMetricPrefix = config.TestRaftMetricPrefix
	c.clusterName = config.ClusterName
	c.config.BindAddress = config.BindAddress
	c.config.HostIPAddress = config.HostIPAddress
	c.config.HTTPPort = config.HTTPPort
	c.config.GRPCPort = config.GRPCPort
	c.config.RaftGRPCPort = config.RaftGRPCPort
	server := &Server{
		Raft: &RaftConfig{
			TimeMultiplier:    1,
			SnapshotInterval:  30 * time.Second,
			SnapshotThreshold: defaultSnapshotThreshold,
		},
	}
	c.config.Server = server
	c.members = config.Members
}

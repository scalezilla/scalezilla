package cluster

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// parseConfig will parse the file configuration for server and client mode
// to validate its properties
func (c *Cluster) parseConfig() error {
	var config Config

	parser := hclparse.NewParser()
	file, diagnostics := parser.ParseHCLFile(c.configFile)
	if diagnostics.HasErrors() {
		return diagnostics.Errs()[0]
	}

	if diagnostics := gohcl.DecodeBody(file.Body, nil, &config); diagnostics.HasErrors() {
		return diagnostics.Errs()[0]
	}

	if config.ClusterName == "" {
		config.ClusterName = defaultClusterName
	}

	if config.DataDir == "" {
		config.DataDir = defaultDataDir
	}

	if config.BindAddress == "" {
		config.BindAddress = defaultBindAddress
	}

	if config.HostIPAddress == "" {
		config.HostIPAddress = defaultHostIPAddress
	}

	if config.HTTPPort == 0 {
		config.HTTPPort = defaultHTTPPort
	}

	if config.GRPCPort == 0 {
		config.GRPCPort = defaultGRPCPort
	}

	if config.RaftGRPCPort == 0 {
		config.RaftGRPCPort = defaultRaftGRPCPort
	}

	if config.Server == nil && config.Client == nil {
		return ErrServerOrClientBlockUndefined
	}

	if config.Server != nil && !config.Server.Enabled && config.Client == nil {
		return ErrServerRaftBlockInvalid
	}

	if config.Server == nil && config.Client != nil && !config.Client.Enabled {
		return ErrClientRaftBlockInvalid
	}

	if config.Server != nil && config.Server.Enabled && config.Server.Raft != nil && config.Server.Raft.BootstrapExpectedSize == 0 || config.Client != nil && config.Client.Enabled && config.Client.Raft != nil && config.Client.Raft.BootstrapExpectedSize == 0 {
		return ErrRaftBootstrapExpectedSizeInvalid
	}

	if config.Server != nil && config.Server.Enabled && config.Server.ClusterJoin == nil {
		return ErrServerClusterJoinBlockInvalid
	}

	if config.Client != nil && config.Client.Enabled && config.Client.ClusterJoin == nil {
		return ErrClientClusterJoinBlockInvalid
	}

	if config.Server != nil && config.Server.ClusterJoin != nil && len(config.Server.ClusterJoin.InitialMembers) == 0 {
		return ErrClusterJoinInitialMembersInvalid
	}

	if config.Client != nil && config.Client.ClusterJoin != nil && len(config.Client.ClusterJoin.InitialMembers) == 0 {
		return ErrClusterJoinInitialMembersInvalid
	}

	if config.Server != nil && config.Server.Enabled {
		c.isVoter = true

		if config.Server.Raft.TimeMultiplier == 0 {
			config.Server.Raft.TimeMultiplier = 1
		}
		if config.Server.Raft.TimeMultiplier > 0 {
			config.Server.Raft.TimeMultiplier = 10
		}

		if config.Server.Raft.SnapshotInterval == 0 {
			config.Server.Raft.SnapshotInterval = defaultSnapshotInterval
		}

		if config.Server.Raft.SnapshotThreshold < 128 {
			config.Server.Raft.SnapshotThreshold = defaultSnapshotThreshold
		}

		if config.Server.ClusterJoin.RetryMax == 0 {
			config.Server.ClusterJoin.RetryMax = defaultClusterJoinRetryMax
		}

		if config.Server.ClusterJoin.RetryInterval == 0 {
			config.Server.ClusterJoin.RetryInterval = defaultClusterJoinRetryInterval
		}
	}

	if config.Client != nil && config.Client.Enabled {
		if config.Client.Raft.TimeMultiplier == 0 {
			config.Client.Raft.TimeMultiplier = 2
		}
		if config.Client.Raft.TimeMultiplier > 0 {
			config.Client.Raft.TimeMultiplier = 10
		}

		if config.Client.Raft.SnapshotInterval == 0 {
			config.Client.Raft.SnapshotInterval = defaultSnapshotInterval
		}

		if config.Client.Raft.SnapshotThreshold < 128 {
			config.Client.Raft.SnapshotThreshold = defaultSnapshotThreshold
		}

		if config.Client.ClusterJoin.RetryMax == 0 {
			config.Client.ClusterJoin.RetryMax = defaultClusterJoinRetryMax
		}

		if config.Client.ClusterJoin.RetryInterval == 0 {
			config.Client.ClusterJoin.RetryInterval = defaultClusterJoinRetryInterval
		}
	}

	c.config = config
	if c.test {
		c.buildDataDir()
	}

	return nil
}

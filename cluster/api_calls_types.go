package cluster

import (
	"context"

	"github.com/rs/zerolog"
)

// ClusterHTTPCallBaseConfig is the configuration used by the cli
// to interact with the cluster
type ClusterHTTPCallBaseConfig struct {
	// Logger is the cluster logger
	Logger *zerolog.Logger

	// HTTPAddress is the address to use to communicate
	// with the cluster api
	HTTPAddress string

	// Context is the context provided by the cli
	// to start the cluster
	Context context.Context

	// OutputFormat can only be table or json
	OutputFormat string
}

// BootstrapClusterHTTPConfig is the configuration used by the cli
// to interact with the cluster
type BootstrapClusterHTTPConfig struct {
	// Token to use to bootstrap the cluster
	Token string

	// Default config
	ClusterHTTPCallBaseConfig
}

// NodesListHTTPConfig is the configuration used by the cli
// to interact with the cluster nodes
type NodesListHTTPConfig struct {
	// Token to use to bootstrap the cluster
	Token string

	// Kind to use to list cluster nodes.
	// Can only be server or client
	Kind string

	// Default config
	ClusterHTTPCallBaseConfig
}

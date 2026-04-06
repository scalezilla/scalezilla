package cluster

import (
	"context"

	"github.com/rs/zerolog"
)

// ClusterHTTPCallBaseConfig is used by the cli to interact with the cluster
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

// BootstrapClusterHTTPConfig is used by the cli to interact with the cluster
type BootstrapClusterHTTPConfig struct {
	// Token to use to bootstrap the cluster
	Token string

	// Default config
	ClusterHTTPCallBaseConfig
}

// NodesListHTTPConfig is used by the cli to interact with the cluster nodes
type NodesListHTTPConfig struct {
	// Token to use to interact with the cluster
	Token string

	// Kind to use to list cluster nodes.
	// Can only be server or client
	Kind string

	// Default config
	ClusterHTTPCallBaseConfig
}

// DeploymentApplyHTTPConfig is used by the cli
// to interact with the server leader node to apply a new deployment
type DeploymentApplyHTTPConfig struct {
	// Token to use to interact with the cluster
	Token string

	// File to use to apply new deployment
	File string

	// Default config
	ClusterHTTPCallBaseConfig

	// osReadFile is a wrapper to osReadFile
	osReadFile func(name string) ([]byte, error)
}

// PodsListHTTPConfig is used by the cli to interact with the cluster nodes
type PodsListHTTPConfig struct {
	// Token to use to interact with the cluster
	Token string

	// Namespace is used to fetch pods for specific or all namespaces
	Namespace string

	// Default config
	ClusterHTTPCallBaseConfig
}

package cluster

import (
	"context"

	"github.com/rs/zerolog"
)

// DeploymentInitialConfig is the configuration used by the cli
// to start a new cluster
type DeploymentInitialConfig struct {
	// Logger is the cluster logger
	Logger *zerolog.Logger

	// ConfigFile is the full path of the config to start the cluster
	ConfigFile string

	// Context is the context provided by the cli
	// to start the cluster
	Context context.Context
}

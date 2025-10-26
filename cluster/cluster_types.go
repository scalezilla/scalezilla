package cluster

import (
	"context"
	"net"
	"os"

	"github.com/rs/zerolog"
)

const scalezillaAppName string = "scalezilla"

// ClusterInitialConfig is the configuration used by the cli
// to start a new cluster
type ClusterInitialConfig struct {
	Logger *zerolog.Logger

	// ConfigPath is the path of the config to start the cluster
	ConfigPath string

	// Dev indicates if we need to start a development cluster
	Dev bool

	// BindAddress is the address to use by the cluster
	BindAddress string

	// HostIPAddress is the ip address to use by the cluster
	HostIPAddress string

	// HTTPPort to use to handle http requests
	HTTPPort uint16

	// RaftGRPCPort is the port for rafty cluster purpose
	RaftGRPCPort uint16

	// GRPCPort is the port for internal cluster purpose
	GRPCPort uint16

	// Members are the rafty cluster members
	Members []string
}

// Cluster holds all required configuration to start the instance
type Cluster struct {
	logger *zerolog.Logger

	// configPath is the path of the config to start the cluster
	configPath string

	// dev indicates if we need to start a development cluster
	dev bool

	// BindAddress is the address to use by the cluster
	bindAddress string

	// HostIPAddress is the ip address to use by the cluster
	hostIPAddress string

	// httpPort to use to handle http requests
	httpPort uint16

	// raftGRPCPort is the port for rafty cluster purpose
	raftGRPCPort uint16

	// grpcPort is the port for internal cluster purpose
	grpcPort uint16

	// Members are the rafty cluster members
	members []string

	// address is the address of the current node
	address net.TCPAddr

	// id is the id of the current node
	id string

	// dataDir is the working directory of this node
	dataDir string

	// quit is the chan used to stop the cluster
	quit chan os.Signal

	// fsm holds requirements to manipulate store
	fsm *fsmState

	// raftyStore holds bold storage config for rafty
	raftyStore raftyStore

	// rafty holds rafty cluster config
	rafty raftyServer

	// apiServer holds the config of the HTTP API server
	apiServer httpServer

	// newRafty is used as a dependency injection
	// newRaftyFunc func() error

	// startRaftyFunc is used as a dependency injection
	startRaftyFunc func() error

	// startAPIServerFunc is used as a dependency injection
	startAPIServerFunc func() error

	// stopAPIServerFunc is used as a dependency injection
	stopAPIServerFunc func() error

	// stopRaftyFunc is used as a dependency injection
	stopRaftyFunc func()

	// raftyStoreCloseFunc is used as a dependency injection
	raftyStoreCloseFunc func() error
}

// httpServer is an interface implements http.Server requirements.
// This will be useful during unit testing
type httpServer interface {
	// ListenAndServe listen and serve the HTTP Server
	ListenAndServe() error

	// Shutdown will shutdown the server
	Shutdown(ctx context.Context) error
}

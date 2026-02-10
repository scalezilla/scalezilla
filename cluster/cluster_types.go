package cluster

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Lord-Y/rafty"
	"github.com/rs/zerolog"
	"github.com/scalezilla/scalezilla/osdiscovery"
	"github.com/scalezilla/scalezilla/scalezillapb"
	"google.golang.org/grpc"
)

const (
	scalezillaAppName  string = "scalezilla"
	defaultClusterName string = "default"
	defaultDataDir     string = "/var/lib/scalezilla"
	defaultNodePool    string = "default"
)

var (
	defaultBindAddress              string        = "127.0.0.1"
	defaultHostIPAddress            string        = "127.0.0.1"
	defaultHTTPPort                 uint16        = 15000
	defaultGRPCPort                 uint16        = 15001
	defaultRaftGRPCPort             uint16        = 15002
	defaultSnapshotInterval         time.Duration = 5 * time.Minute
	defaultSnapshotThreshold        uint64        = 128
	defaultClusterJoinRetryInterval time.Duration = 15 * time.Second
	defaultClusterJoinRetryMax      uint16        = 5
	defaultMaxAppendEntries         uint64        = 1024
	defaultGrpcForceTimeout         time.Duration = 30 * time.Second
)

// ClusterInitialConfig is the configuration used by the cli
// to start a new cluster
type ClusterInitialConfig struct {
	// Logger is the cluster logger
	Logger *zerolog.Logger

	// ConfigFile is the full path of the config to start the cluster
	ConfigFile string

	// Dev indicates if we need to start a development cluster
	Dev bool

	// Test indicates if we need to override some settings
	// for unit testing
	Test bool

	// BindAddress is the address to use by the cluster
	BindAddress string

	// HostIPAddress is the ip address to use by the cluster
	HostIPAddress string

	// HTTPPort to use to handle http requests
	HTTPPort uint16

	// GRPCPort is the port for internal cluster purpose
	GRPCPort uint16

	// RaftGRPCPort is the port for rafty cluster purpose
	RaftGRPCPort uint16

	// Members are the rafty cluster members
	Members []string

	// ClusterName is the name of the current cluster
	ClusterName string

	// NodePool is the pool associated to the current node
	NodePool string

	// TestRaftMetricPrefix is only used during unit testing to override
	// raft metric prefix
	TestRaftMetricPrefix string

	// Context is the context provided by the cli
	// to start the cluster
	Context context.Context
}

// Cluster holds all required configuration to start the instance
type Cluster struct {
	logger *zerolog.Logger

	// mu is used to ensure lock concurrency
	mu sync.Mutex

	// wg is the goroutines tracker
	wg sync.WaitGroup

	// configFile is the full path of the config to start the cluster
	configFile string

	// clusterName is the name of the current cluster
	clusterName string

	// dev indicates if we need to start a development cluster
	dev bool

	// test indicates if we need to override some settings
	// for unit testing
	test bool

	// isVoter statuates if the current node is a voting member node
	isVoter bool

	// Members are the rafty cluster members
	members []string

	// grpcAddress is the address of the current node
	grpcAddress net.TCPAddr

	// raftyAddress is the address of the current node
	raftyAddress net.TCPAddr

	// id is the id of the current node
	id string

	// nodePool is the pool associated to the current node
	nodePool string

	// ctx is the context to use to stop the cluster
	ctx context.Context

	// fsm holds requirements to manipulate store
	fsm *fsmState

	// raftyStore holds bold storage config for rafty
	raftyStore raftyStore

	// rafty holds rafty cluster config
	rafty raftyServer

	// apiServer holds the config of the HTTP API server
	apiServer httpServer

	// di holds all dependency injection funcs
	di dependencyInjections

	// grpcListen holds the listener of the grpc server
	grpcListen net.Listener

	// grpcServer holds the grpc server
	grpcServer *grpc.Server

	// scalezillapb.ScalezillaServer holds the interface
	// to interact with the grpc server
	scalezillapb.ScalezillaServer

	// grpcForceTimeout is used to force stop the grpc server
	grpcForceTimeout time.Duration

	// raftMetricPrefix is used to prefix rafty metrics
	// with the provided value
	raftMetricPrefix string

	// config is the configuration to use to start the cluster
	config Config

	// systemInfo holds os discovery requirements
	systemInfo *osdiscovery.SystemInfo

	// isRunning is a helper indicating is the node is up or down.
	// It set to false, it will reject all incoming grpc requests
	// with shutting down error
	isRunning atomic.Bool

	// bootstrapExpectedSizeReach is an atomic bool flag set
	// to check if the bootstrap size is reached
	bootstrapExpectedSizeReach atomic.Bool

	// bootstrapExpectedSize is a counter that goes
	// with bootstrapExpectedSizeReach variable
	bootstrapExpectedSize atomic.Uint64

	// clientContactedServer is an atomic bool flag set
	// to check if the client has successfully contacted
	// the server(s). It will then ask to be part of the cluster
	// once boostrapped
	clientContactedServer atomic.Bool

	// nodeMap is a map of all nodes in the cluster
	nodeMap map[string]*nodeMap

	// nodeMapMu is used to ensure lock concurrency
	nodeMapMu sync.Mutex

	// connectionManager holds connections for all nodes
	connectionManager connectionManager

	// rpcServicePortsDiscoveryChanReq is used by the grpc
	// receiver to respond back to the caller
	rpcServicePortsDiscoveryChanReq chan RPCRequest

	// rpcServicePortsDiscoveryChanResp is used to receive
	// answers sent by actual node
	rpcServicePortsDiscoveryChanResp chan RPCResponse

	// rpcServiceNodePollingChanReq is used by the grpc
	// receiver to respond back to the caller
	rpcServiceNodePollingChanReq chan RPCRequest

	// rpcServiceNodePollingChanResp is used to receive
	// answers sent by actual node
	rpcServiceNodePollingChanResp chan RPCResponse

	// checkBootstrapSizeDuration is the frequency at which
	// to make rpc calls to other nodes to satisfy
	// bootstrapExpectedSize variable
	checkBootstrapSizeDuration time.Duration

	// nodePollingTimer is the frequency at which the node
	// will send node polling rpc request to other nodes
	nodePollingTimer time.Duration
}

// dependencyInjections is a struct holding all
// dependency injection funcs
type dependencyInjections struct {
	// newRaftyFunc is used as a dependency injection
	newRaftyFunc func() (*rafty.Rafty, error)

	// startRaftyFunc is used as a dependency injection
	startRaftyFunc func() error

	// startAPIServerFunc is used as a dependency injection
	startAPIServerFunc func() error

	// stopAPIServerFunc is used as a dependency injection
	stopAPIServerFunc func() error

	// startGRPCServerFunc is used as a dependency injection
	startGRPCServerFunc func() error

	// grpcListenFunc holds the listener of the grpc server
	grpcListenFunc func(network string, address string) (net.Listener, error)

	// newGRPCServerFunc is used as a dependency injection
	newGRPCServerFunc func(opt ...grpc.ServerOption) *grpc.Server

	// stopGRPCServerFunc is used as a dependency injection
	stopGRPCServerFunc func()

	// stopRaftyFunc is used as a dependency injection
	stopRaftyFunc func()

	// raftyStoreCloseFunc is used as a dependency injection
	raftyStoreCloseFunc func() error

	// checkSystemInfoFunc is used as a dependency injection
	checkSystemInfoFunc func() error

	// osdiscoveryFunc is used as a dependency injection
	osdiscoveryFunc func() *osdiscovery.SystemInfo

	// checkBootstrapSizeFunc is used as a dependency injection
	checkBootstrapSizeFunc func()

	// sendRPCFunc is used as a dependency injection
	sendRPCFunc func(address string, client scalezillapb.ScalezillaClient, request RPCRequest)

	// aclTokenEncodeCommandFunc is used as a dependency injection
	aclTokenEncodeCommandFunc func(cmd aclTokenCommand, w io.Writer) error
}

// httpServer is an interface implements http.Server requirements.
// This will be useful during unit testing
type httpServer interface {
	// ListenAndServe listen and serve the HTTP Server
	ListenAndServe() error

	// Shutdown will shutdown the server
	Shutdown(ctx context.Context) error
}

// Config is the configuration to use to start
// the cluster
type Config struct {
	// Hostname is the name of the current host
	Hostname string `hcl:"hostname,optional"`

	// ClusterName is the name of the current cluster
	ClusterName string `hcl:"cluster_name,optional"`

	// DataDir is where node data will be stored
	DataDir string `hcl:"data_dir,optional"`

	// BindAddress is the address to use by the cluster
	BindAddress string `hcl:"bind_address,optional"`

	// HostIPAddress is the ip address to use by the cluster
	HostIPAddress string `hcl:"host_ip_address,optional"`

	// HTTPPort to use to handle http requests
	HTTPPort uint16 `hcl:"http_port,optional"`

	// RaftGRPCPort is the port for rafty cluster purpose
	RaftGRPCPort uint16 `hcl:"raft_grpc_port,optional"`

	// GRPCPort is the port for internal cluster purpose
	GRPCPort uint16 `hcl:"grpc_port,optional"`

	// Metadata is the ip address to use by the cluster
	Metadata map[string]string `hcl:"metadata,optional"`

	// Server holds controle plane config
	Server *Server `hcl:"server,block"`

	// Client holds data plane config
	Client *Client `hcl:"client,block"`
}

// Server holds the requirements to start current node
// as a cluster
type Server struct {
	// Enabled when set to true indicates we are in server mode
	Enabled bool `hcl:"enabled"`

	// Raft holds config related to raft consensus protocol
	Raft *RaftConfig `hcl:"raft,block"`

	// ClusterJoin holds requirements to join the cluster
	ClusterJoin *ClusterJoin `hcl:"cluster_join,block"`
}

// Server holds the requirements to start current node
// as a cluster
type Client struct {
	// Enabled when set to true indicates we are in server mode
	Enabled bool `hcl:"enabled"`

	// Raft holds config related to raft consensus protocol
	Raft *RaftConfig `hcl:"raft,block"`

	// ClusterJoin holds requirements to join the cluster
	ClusterJoin *ClusterJoin `hcl:"cluster_join,block"`

	// NodePool is the pool associated to the current node
	NodePool *string `hcl:"node_pool"`
}

// RaftConfig holds the requirements to start the raft cluster
type RaftConfig struct {
	// // dataDir is where node data will be stored
	// dataDir string `hcl:"data_dir,optional"`

	// // logger expose zerolog config so it can be override
	// logger *zerolog.Logger

	// BootstrapExpectedSize is the number of node to wait for
	// before bootstrapping the cluster
	BootstrapExpectedSize uint64 `hcl:"bootstrap_expected_size"`

	// TimeMultiplier is a scaling factor that will be used during election timeout and heartbeats checks.
	// Default to 1 for server mode.
	// Default to 2 for client mode.
	// Max is 10.
	TimeMultiplier uint `hcl:"time_multiplier,optional"`

	// SnapshotInterval is the interval at which a snapshot will be taken. It will be randomize with this minimum value in order to prevent all nodes to take a snapshot at the same time.
	// Default to 5 minutes
	SnapshotInterval time.Duration `hcl:"snapshot_interval,optional"`

	// SnapshotThreshold is the threshold that must be reached before taking a snapshot.
	// It prevent to take snapshots to frequently.
	// Default to 128
	SnapshotThreshold uint64 `hcl:"snapshot_threshold,optional"`
}

// ClusterJoin holds requirements to join the cluster
type ClusterJoin struct {
	// InitialMembers is the list of members to contact
	// to join the cluster.
	// Format is [ "x.x.x.x", "y.y.y.y", "z.z.z.z:15002"]
	// When port is not specified, it defaults to 15002
	InitialMembers []string `hcl:"initial_members"`

	// RetryMax is the maximum retry to contact initial members.
	// Default to 5
	RetryMax uint16 `hcl:"retry_max,optional"`

	// RetryInterval is a timeout at which the current node will
	// try to contact the initial members.
	// Default to 15s
	RetryInterval time.Duration `hcl:"retry_interval,optional"`
}

// nodeMap is a map of all nodes in the cluster
type nodeMap struct {
	// IsVoter when set to true means it's a server node
	IsVoter bool

	// ID is the ID of the node
	ID string

	// Address is the host ip of the node
	Address string

	// HTTPPort is the http port of the node
	HTTPPort uint32

	// GRPCPort is the http port of the node
	GRPCPort uint32

	// RaftyPort is the raft port of the node
	RaftyPort uint32

	// NodePool is the node pool of the node
	NodePool string

	// SystemInfo holds os discovery requirements
	SystemInfo osdiscovery.SystemInfo

	// Metadata holds node metadata
	Metadata map[string]string
}

// connectionManager is used to manage all grpc connections nodes
type connectionManager struct {
	// mu is used to ensure lock concurrency
	mu sync.Mutex

	// connections holds gprc server connection for all clients
	connections map[string]*grpc.ClientConn

	// clients holds gprc rafty client for all clients
	clients map[string]scalezillapb.ScalezillaClient
}

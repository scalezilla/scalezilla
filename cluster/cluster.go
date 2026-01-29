package cluster

import (
	"net"
	"time"

	"github.com/scalezilla/scalezilla/osdiscovery"
	"github.com/scalezilla/scalezilla/scalezillapb"
	"google.golang.org/grpc"
)

// NewCluster build the requirements to start the cluster
func NewCluster(config ClusterInitialConfig) (*Cluster, error) {
	c := &Cluster{
		logger:           config.Logger,
		configFile:       config.ConfigFile,
		ctx:              config.Context,
		test:             config.Test,
		grpcForceTimeout: defaultGrpcForceTimeout,
		nodePool:         config.NodePool,
		nodeMap:          make(map[string]*nodeMap),
		connectionManager: connectionManager{
			connections: make(map[string]*grpc.ClientConn),
			clients:     make(map[string]scalezillapb.ScalezillaClient),
		},
		rpcServicePortsDiscoveryChanReq:  make(chan RPCRequest),
		rpcServicePortsDiscoveryChanResp: make(chan RPCResponse),
		checkBootstrapSizeDuration:       5 * time.Second,
	}

	c.di.newRaftyFunc = c.newRafty
	c.di.startRaftyFunc = c.startRafty
	c.di.startAPIServerFunc = c.startAPIServer
	c.di.stopAPIServerFunc = c.stopAPIServer
	c.di.startGRPCServerFunc = c.startGRPCServer
	c.di.grpcListenFunc = net.Listen
	c.di.stopGRPCServerFunc = c.stopGRPCServer
	c.di.stopRaftyFunc = c.stopRafty
	c.di.raftyStoreCloseFunc = c.raftyStoreClose
	c.di.checkSystemInfoFunc = c.checkSystemInfo
	c.di.osdiscoveryFunc = osdiscovery.NewSystemInfo
	c.di.checkBootstrapSizeFunc = c.checkBootstrapSize
	c.di.sendRPCFunc = c.sendRPC
	c.raftMetricPrefix = scalezillaAppName
	c.di.aclTokenEncodeCommandFunc = aclTokenEncodeCommand
	c.di.sendRPCFunc = c.sendRPC

	c.buildDataDir()
	if config.Dev {
		c.buildDevConfig(config)
		return c, nil
	}

	if err := c.parseConfig(); err != nil {
		return nil, err
	}

	return c, nil
}

// Start will start the cluster
func (c *Cluster) Start() error {
	c.buildAddressAndID()

	var err error
	if err := c.di.checkSystemInfoFunc(); err != nil {
		return err
	}

	if err := c.di.startGRPCServerFunc(); err != nil {
		return err
	}

	if c.rafty, err = c.di.newRaftyFunc(); err != nil {
		return err
	}

	if err := c.di.startRaftyFunc(); err != nil {
		return err
	}

	c.newAPIServer()
	if err := c.di.startAPIServerFunc(); err != nil {
		return err
	}
	c.logger.Info().Msg("server started successfully")

	c.isRunning.Store(true)
	<-c.ctx.Done()

	c.isRunning.Store(false)
	if err := c.di.stopAPIServerFunc(); err != nil {
		return err
	}

	c.di.stopGRPCServerFunc()
	c.di.stopRaftyFunc()

	if err := c.di.raftyStoreCloseFunc(); err != nil {
		return err
	}
	c.wg.Wait()

	c.logger.Info().Msg("server stopped successfully")
	return nil
}

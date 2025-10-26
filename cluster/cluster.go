package cluster

import "fmt"

// NewCluster build the requirements to start the cluster
func NewCluster(config ClusterInitialConfig) *Cluster {
	c := &Cluster{
		logger:        config.Logger,
		configPath:    config.ConfigPath,
		dev:           config.Dev,
		bindAddress:   config.BindAddress,
		hostIPAddress: config.HostIPAddress,
		httpPort:      config.HTTPPort,
		raftGRPCPort:  config.GRPCPort,
		grpcPort:      config.GRPCPort,
		members:       config.Members,
	}

	c.startRaftyFunc = c.startRafty
	c.startAPIServerFunc = c.startAPIServer
	c.stopAPIServerFunc = c.stopAPIServer
	c.stopRaftyFunc = c.stopRafty
	c.raftyStoreCloseFunc = c.raftyStoreClose
	return c
}

// Start will start the cluster
func (c *Cluster) Start() error {
	c.buildAddressAndID()
	c.buildDataDir()
	c.buildSignal()

	var err error
	if c.rafty, err = c.newRafty(scalezillaAppName); err != nil {
		return err
	}
	c.newAPIServer()

	if err := c.startRaftyFunc(); err != nil {
		return err
	}

	if err := c.startAPIServerFunc(); err != nil {
		return err
	}
	c.logger.Info().Msg("server started successfully")

	<-c.quit

	fmt.Println("FUCK")

	if err := c.stopAPIServerFunc(); err != nil {
		return err
	}

	c.stopRaftyFunc()
	c.logger.Info().Msg("server stopped successfully")

	if err := c.raftyStoreCloseFunc(); err != nil {
		return err
	}
	return nil
}

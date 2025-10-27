package cluster

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

	c.newRaftyFunc = c.newRafty
	c.startRaftyFunc = c.startRafty
	c.startAPIServerFunc = c.startAPIServer
	c.stopAPIServerFunc = c.stopAPIServer
	c.stopRaftyFunc = c.stopRafty
	c.raftyStoreCloseFunc = c.raftyStoreClose

	c.buildDataDir()
	c.buildSignal()
	return c
}

// Start will start the cluster
func (c *Cluster) Start() error {
	c.buildAddressAndID()

	var err error
	if c.rafty, err = c.newRaftyFunc(scalezillaAppName); err != nil {
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

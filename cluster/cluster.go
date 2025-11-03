package cluster

import "github.com/scalezilla/scalezilla/osdiscovery"

// NewCluster build the requirements to start the cluster
func NewCluster(config ClusterInitialConfig) (*Cluster, error) {
	c := &Cluster{
		logger:     config.Logger,
		configFile: config.ConfigFile,
		ctx:        config.Context,
		test:       config.Test,
	}

	c.newRaftyFunc = c.newRafty
	c.startRaftyFunc = c.startRafty
	c.startAPIServerFunc = c.startAPIServer
	c.stopAPIServerFunc = c.stopAPIServer
	c.stopRaftyFunc = c.stopRafty
	c.raftyStoreCloseFunc = c.raftyStoreClose
	c.raftMetricPrefix = scalezillaAppName
	c.checkSystemInfoFunc = c.checkSystemInfo
	c.osdiscoveryFunc = osdiscovery.NewSystemInfo

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
	if err := c.checkSystemInfoFunc(); err != nil {
		return err
	}

	if c.rafty, err = c.newRaftyFunc(); err != nil {
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

	<-c.ctx.Done()

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

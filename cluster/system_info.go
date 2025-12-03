package cluster

// checkSystemInfo will check if all requirements are met
// to start the cluster
func (c *Cluster) checkSystemInfo() error {
	s := c.di.osdiscoveryFunc()
	if s.Cgroups.Version != 2 {
		return ErrCgroupV2Required
	}
	c.systemInfo = s

	return nil
}

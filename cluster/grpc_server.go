package cluster

import (
	"time"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// startGRPCServer will start the internal grpc server
func (c *Cluster) startGRPCServer() error {
	var err error
	if c.grpcListen, err = c.di.grpcListenFunc(c.grpcAddress.Network(), c.grpcAddress.String()); err != nil {
		return err
	}

	newServer := c.di.newGRPCServerFunc
	if newServer == nil {
		newServer = grpc.NewServer
	}

	c.mu.Lock()
	c.grpcServer = newServer()
	c.mu.Unlock()
	scalezillapb.RegisterScalezillaServer(c.grpcServer, c.ScalezillaServer)

	errChan := make(chan error, 1)
	go func() {
		errChan <- c.grpcServer.Serve(c.grpcListen)
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(10 * time.Millisecond):
	}

	c.wg.Go(c.grpcLoop)
	c.wg.Go(c.checkBootstrapSize)
	return nil
}

// stopGRPCServer will stop the internal grpc server
func (c *Cluster) stopGRPCServer() {
	c.mu.Lock()
	c.grpcServer.GracefulStop()
	c.mu.Unlock()
}

// getClient return rpc connection client
func (c *Cluster) getClient(address string) scalezillapb.ScalezillaClient {
	c.connectionManager.mu.Lock()
	defer c.connectionManager.mu.Unlock()

	if client, ok := c.connectionManager.clients[address]; ok {
		return client
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// in our case for now, err is always nil so no need to check it
	conn, _ := grpc.NewClient(
		address,
		opts...,
	)

	c.connectionManager.connections[address] = conn
	c.connectionManager.clients[address] = scalezillapb.NewScalezillaClient(conn)
	return c.connectionManager.clients[address]
}

// checkBootstrapSize will wait to get the expected bootstrap
// size by sending ServicePortsDiscovery requests
func (c *Cluster) checkBootstrapSize() {
	if c.dev || !c.isVoter {
		return
	}

	timer := time.NewTicker(c.checkBootstrapSizeDuration)
	defer timer.Stop()

	go c.reqServicePortsDiscovery()
	for !c.bootstrapExpectedSizeReach.Load() {
		select {
		case <-c.ctx.Done():
			return

		case <-timer.C:
			if c.bootstrapExpectedSize.Load()+1 == c.config.Server.Raft.BootstrapExpectedSize {
				c.bootstrapExpectedSizeReach.Store(true)
				return
			}
			go c.reqServicePortsDiscovery()
		}
	}
}

package cluster

import (
	"time"

	"github.com/scalezilla/scalezilla/scalezillapb"
	"google.golang.org/grpc"
)

// startGRPCServer will start the internal grpc server
func (c *Cluster) startGRPCServer() error {
	var err error
	if c.grpcListen, err = c.grpcListenFunc(c.grpcAddress.Network(), c.grpcAddress.String()); err != nil {
		return err
	}

	newServer := c.newGRPCServerFunc
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

	return nil
}

// stopGRPCServer will stop the internal grpc server
func (c *Cluster) stopGRPCServer() {
	c.mu.Lock()
	c.grpcServer.GracefulStop()
	c.mu.Unlock()
}

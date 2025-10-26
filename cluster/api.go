package cluster

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// newAPIServer will build the api server config
func (c *Cluster) newAPIServer() {
	c.apiServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", c.httpPort),
		Handler: c.newApiRouters(),
	}
}

// startAPI will start the api server
func (c *Cluster) startAPIServer() error {
	errChan := make(chan error, 1)
	defer close(errChan)
	go func() {
		if err := c.apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(10 * time.Millisecond):
		return nil
	}
}

// stopAPIServer will stop the api server
func (c *Cluster) stopAPIServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := c.apiServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

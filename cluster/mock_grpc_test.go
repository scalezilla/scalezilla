package cluster

import (
	"fmt"
	"net"
)

type badListener struct{}

func (badListener) Accept() (net.Conn, error) {
	return nil, fmt.Errorf("accept failed")
}

func (badListener) Close() error {
	return nil
}

func (badListener) Addr() net.Addr {
	return &net.TCPAddr{}
}

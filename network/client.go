package network

import (
	"net"
)

func Connect(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Conn) StoreBlob() {
}

func (c *Conn) Close() {
}

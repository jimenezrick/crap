package network

import (
	"bufio"
	"io"
	"net"
)

import "crap/kvmap"

type Conn struct {
	sock net.Conn
	io.ReadWriter
}

func newConn(sock net.Conn) *Conn {
	return &Conn{
		sock,
		bufio.NewReadWriter(
			bufio.NewReader(sock),
			bufio.NewWriter(sock)),
	}
}

func Connect(addr string) (*Conn, error) {
	sock, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return newConn(sock), nil
}

func (c *Conn) StoreBlob(blob io.Reader, size uint32) error {
	req := kvmap.NewWith("request", "store")
	if err := c.WriteJSONFrame(req); err != nil {
		return err
	}

	if err := c.WriteBlobFrameFrom(blob, size); err != nil {
		return err
	}

	req = kvmap.NewWith("key", "bogus") // TODO
	if err := c.WriteJSONFrame(req); err != nil {
		return err
	}

	// XXX XXX XX
	var res result
	if err := c.ReadJSONFrame(&res); err != nil {
		return err
	}
	println("RESULT: ", res.val, res.info)
	// XXX XXX XXX

	return nil
}

func (c *Conn) Flush() error {
	return c.ReadWriter.(*bufio.ReadWriter).Flush()
}

func (c *Conn) Close() error {
	return c.sock.Close()
}

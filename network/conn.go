package network

import (
	"os"
	"io"
	"bufio"
	"net"
)

type conn struct {
	sock net.Conn
	io.ReadWriter
}

type request struct {
	Req string "request"
	Key string "key"
}

type result struct {
	Res  string "result"
	Info string "info"
}

func newConn(sock net.Conn) *conn {
	return &conn{
		sock,
		bufio.NewReadWriter(bufio.NewReader(sock), bufio.NewWriter(sock)),
	}
}

func Connect(addr string) (*conn, error) {
	sock, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return newConn(sock), nil
}

func (c *conn) StoreBlob(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	req := request{"store", name} // XXX
	size := uint32(info.Size())   // XXX
	if err = c.WriteJSONFrame(req); err != nil {
		return err
	}
	if err = c.WriteBlobFrameFrom(bufio.NewReader(file), size); err != nil {
		return err
	}

	return nil
}

func (c *conn) Flush() error {
	return c.ReadWriter.(*bufio.ReadWriter).Flush()
}

func (c *conn) Close() error {
	return c.sock.Close()
}

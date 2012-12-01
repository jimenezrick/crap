package network

import (
	"bufio"
	"os"
	"io"
	"net"
)

import (
	"crap/hashed"
	"crap/util"
)

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

func (c *Conn) StoreBlob(file *os.File) (string, error) {
	reader := hashed.NewSHA1FileReader(file)

	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	if err := c.WriteJSONFrame(request{"store"}); err != nil {
		return "", err
	}

	if err := c.WriteBlobFrameFrom(reader, uint32(info.Size())); err != nil {
		return "", err
	}

	key := util.HexHash(reader)
	if err := c.WriteJSONFrame(keyRequest{key}); err != nil {
		return "", err
	}

	var res response
	if err := c.ReadJSONFrame(&res); err != nil {
		return "", err
	}

	if res.Val != "ok" {
		return "", responseError(res)
	}

	return key, nil
}

func (c *Conn) Flush() error {
	return c.ReadWriter.(*bufio.ReadWriter).Flush()
}

func (c *Conn) Close() error {
	return c.sock.Close()
}

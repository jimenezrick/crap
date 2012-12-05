package network

import (
	"bufio"
	"io"
	"net"
	"os"
)

import (
	"crap/hashed"
	"crap/store"
	"crap/util"
)

type Conn struct {
	store *store.Store
	sock  net.Conn
	io.ReadWriter
}

func newConn(store *store.Store, sock net.Conn) *Conn {
	return &Conn{
		store,
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
	return newConn(nil, sock), nil
}

func (c *Conn) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := c.Read(buf)
	if err != nil {
		return 0, err
	}

	return byte(buf[0]), nil
}

func (c *Conn) StoreBlob(file *os.File) (string, error) {
	reader := hashed.NewSha1FileReader(file)

	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	size := uint64(info.Size())

	if err := c.WriteJSONFrame(request{"store", "after", size, false}); err != nil {
		return "", err
	}

	if err := c.WriteBlobFrameFrom(reader, size); err != nil {
		return "", err
	}

	key := util.HexHash(reader)
	if err := c.WriteJSONFrame(request{Key: key}); err != nil {
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

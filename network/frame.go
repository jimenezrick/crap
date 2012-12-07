package network

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

import "crap/log"

const maxJSONSize = 4096

var errSize error = errors.New("network: frame too big")
var errMoreData error = errors.New("network: expecting more frame data")

func (c *Conn) readFrameSize() (int64, error) {
	size, err := binary.ReadUvarint(c)
	return int64(size), err
}

func (c *Conn) writeFrameSize(size int64) error {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(size))
	buf = buf[:n]

	_, err := c.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) readFrameBody(max int64) ([]byte, error) {
	size, err := c.readFrameSize()
	if err != nil {
		return nil, err
	}

	if size > max {
		return nil, errSize
	}

	buf := make([]byte, size)
	n, err := io.ReadFull(c, buf)
	if int64(n) != size {
		return nil, errMoreData
	} else if err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *Conn) writeFrameBody(body []byte) error {
	if err := c.writeFrameSize(int64(len(body))); err != nil {
		return err
	}

	_, err := c.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) ReadJSONFrame(obj interface{}) error {
	body, err := c.readFrameBody(maxJSONSize)
	if err != nil {
		return err
	}
	log.Debug.Printf("JSON frame received: %s", body)

	if err = json.Unmarshal(body, obj); err != nil {
		return err
	}

	return nil
}

func (c *Conn) WriteJSONFrame(obj interface{}) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = c.writeFrameBody(body)
	if err != nil {
		return err
	}
	log.Debug.Printf("JSON frame sent: %s", body)

	return c.Flush()
}

func (c *Conn) ReadBlobFrameTo(to io.Writer) error {
	size, err := c.readFrameSize()
	if err != nil {
		return err
	}

	n, err := io.CopyN(to, c, size)
	if n != size {
		return errMoreData
	} else if err != nil {
		return err
	}
	log.Debug.Printf("%d bytes blob frame received", size)

	return nil
}

func (c *Conn) WriteBlobFrameFrom(from io.Reader, size int64) error {
	err := c.writeFrameSize(size)
	if err != nil {
		return err
	}

	_, err = io.CopyN(c, from, size)
	if err != nil {
		return err
	}
	log.Debug.Printf("%d bytes blob frame sent", size)

	return c.Flush()
}

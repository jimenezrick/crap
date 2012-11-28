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

func (c *Conn) readFrameSize() (size uint32, err error) {
	err = binary.Read(c, binary.BigEndian, &size)
	return
}

func (c *Conn) writeFrameSize(size uint32) error {
	return binary.Write(c, binary.BigEndian, size)
}

func (c *Conn) readFrameBody(max uint32) ([]byte, error) {
	size, err := c.readFrameSize()
	if err != nil {
		return nil, err
	}

	if size > max {
		return nil, errSize
	}

	buf := make([]byte, size)
	n, err := io.ReadFull(c, buf)
	if uint32(n) != size {
		return nil, errMoreData
	} else if err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *Conn) writeFrameBody(body []byte) error {
	if err := c.writeFrameSize(uint32(len(body))); err != nil {
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

	n, err := io.CopyN(to, c, int64(size))
	if n != int64(size) {
		return errMoreData
	} else if err != nil {
		return err
	}
	log.Debug.Printf("%d bytes blob frame received", size)

	return nil
}

func (c *Conn) WriteBlobFrameFrom(from io.Reader, size uint32) error {
	err := c.writeFrameSize(size)
	if err != nil {
		return err
	}

	_, err = io.CopyN(c, from, int64(size))
	if err != nil {
		return err
	}
	log.Debug.Printf("%d bytes blob frame sent", size)

	return c.Flush()
}

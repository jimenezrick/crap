package network

import (
	"io"
	"errors"
	"encoding/binary"
	"encoding/json"
)

import "crap/config"

var errJSONSize error = errors.New("crap/net: JSON frame too big")
var errMoreData error = errors.New("crap/net: expecting more frame data")

func (c conn) readFrameSize() (size uint32, err error) {
	err = binary.Read(c, binary.BigEndian, &size)
	return
}

func (c conn) writeFrameSize(size uint32) error {
	return binary.Write(c, binary.BigEndian, size)
}

func (c conn) readFrameBody() ([]byte, error) {
	size, err := c.readFrameSize()
	if err != nil {
		return nil, err
	}

	if size > uint32(config.GetInt("network.max_json_frame_size")) {
		return nil, errJSONSize
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

func (c conn) writeFrameBody(body []byte) error {
	if err := c.writeFrameSize(uint32(len(body))); err != nil {
		return err
	}

	_, err := c.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (c conn) ReadJSONFrame(obj interface{}) error {
	buf, err := c.readFrameBody()
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf, obj); err != nil {
		return err
	}

	return nil
}

func (c conn) WriteJSONFrame(obj interface{}) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = c.writeFrameBody(body)
	if err != nil {
		return err
	}

	return c.Flush()
}

func (c conn) ReadBlobFrameTo(to io.Writer) error {
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

	return nil
}

func (c conn) WriteBlobFrameFrom(from io.Reader, size uint32) error {
	err := c.writeFrameSize(size)
	if err != nil {
		return err
	}

	_, err = io.CopyN(c, from, int64(size))
	if err != nil {
		return err
	}

	return c.Flush()
}

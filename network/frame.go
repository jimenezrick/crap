package network

// XXX: CopyBLOBFrame(r io.Reader, w io.Writer): for copying the blob in the temp file

import (
	"io"
	"errors"
	"encoding/binary"
	"encoding/json"
)

import "crap/config"

func readFrameSize(r io.Reader) (size uint32, err error) {
	err = binary.Read(r, binary.BigEndian, &size)
	return
}

func readFrameBody(r io.Reader) ([]byte, error) {
	size, err := readFrameSize(r)
	if err != nil {
		return nil, err
	}

	if int(size) > config.GetInt("net.max_json_frame_size") {
		return nil, errors.New("crap/net: JSON frame too big")
	}

	buf := make([]byte, size)

	n, err := io.ReadFull(r, buf)
	if n != int(size) {
		return nil, errors.New("crap/net: expecting more frame data")
	} else if err != nil {
		return nil, err
	}

	return buf, nil
}

func ReadJSONFrame(r io.Reader, obj interface{}) error {
	buf, err := readFrameBody(r)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf, obj); err != nil {
		return err
	}

	return nil
}

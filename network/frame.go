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

func readFrameSize(r io.Reader) (size uint32, err error) {
	err = binary.Read(r, binary.BigEndian, &size)
	return
}

func writeFrameSize(w io.Writer, size uint32) error {
	return binary.Write(w, binary.BigEndian, size)
}

func readFrameBody(r io.Reader) ([]byte, error) {
	size, err := readFrameSize(r)
	if err != nil {
		return nil, err
	}

	if size > uint32(config.GetInt("network.max_json_frame_size")) {
		return nil, errJSONSize
	}

	buf := make([]byte, size)

	n, err := io.ReadFull(r, buf)
	if uint32(n) != size {
		return nil, errMoreData
	} else if err != nil {
		return nil, err
	}

	return buf, nil
}

func writeFrameBody(w io.Writer, body []byte) error {
	if err := writeFrameSize(w, uint32(len(body))); err != nil {
		return err
	}

	_, err := w.Write(body)
	if err != nil {
		return err
	}

	return nil
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

func WriteJSONFrame(w io.Writer, obj interface{}) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = writeFrameBody(w, body)
	if err != nil {
		return err
	}

	return nil
}

func ReadBlobFrameTo(r io.Reader, w io.Writer) error {
	size, err := readFrameSize(r)
	if err != nil {
		return err
	}

	n, err := io.CopyN(w, r, int64(size))
	if n != int64(size) {
		return errMoreData
	} else if err != nil {
		return err
	}

	return nil
}

func WriteBlobFrameFrom(w io.Writer, r io.Reader, size uint32) error {
	err := writeFrameSize(w, size)
	if err != nil {
		return err
	}

	_, err = io.CopyN(w, r, int64(size))
	if err != nil {
		return err
	}

	return nil
}

package network

import (
	"errors"
	"net"
)

// XXX: See http://code.google.com/p/go/issues/detail?id=4373
var ErrClosing = errors.New("use of closed network connection")

func IsErrClosing(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err.Error() == ErrClosing.Error()
	}
	return false
}

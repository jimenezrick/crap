package network

import (
	"errors"
	"net"
)

// XXX: See http://code.google.com/p/go/issues/detail?id=4373
var errClosing = errors.New("use of closed network connection")

func IsClosing(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err.Error() == errClosing.Error()
	}
	return false
}

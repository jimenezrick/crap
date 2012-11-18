package service

import (
	"net"
	"time"
	"bytes"
)

const (
	// XXX: CamelCase
	service_address = "239.0.0.1:8000"
	service_timeout = time.Millisecond * 100
	service_retries = 3
)

const closed_connection_error = "use of closed network connection"

func Discover(name string) (nodes []string, err error) {
	gaddr, err := net.ResolveUDPAddr("udp", service_address)
	if err != nil {
		return
	}

	conn, err := net.ListenMulticastUDP("udp", nil, gaddr)
	if err != nil {
		return
	}
	defer conn.Close()

	deadline := time.Now().Add(service_timeout)
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		return
	}

	go func() {
		for {
			_, err := conn.WriteTo([]byte("service_request:" + name), gaddr)
			if err != nil && err.Error() == closed_connection_error {
				return
			} else if err != nil {
				panic(err)
			}
			time.Sleep(service_timeout / service_retries)
		}
	}()

	nodes = make([]string, 0, 16)
	responses := make(map[string] bool)
	buf := make([]byte, 2048)

	for {
		n, addr, uerr := conn.ReadFromUDP(buf)
		if nerr, ok := uerr.(net.Error); ok && nerr.Timeout() {
			for ip := range responses {
				nodes = append(nodes, ip)
			}
			return
		} else if uerr != nil {
			return nil, uerr
		}

		if bytes.Equal(buf[:n], []byte("service_offer:" + name)) {
			responses[addr.IP.String()] = true
		}
	}

	return
}

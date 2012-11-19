package service

import (
	"net"
	"strings"
)

const (
	opQuit = iota
	opRequest
	opRegister
	opUnregister
)

type cmd struct {
	op   int
	serv string
}

type Responder struct {
	gaddr   *net.UDPAddr
	conn    *net.UDPConn
	cmdChan chan cmd
	table   map[string]bool
}

func (res *Responder) Start() error {
	gaddr, err := net.ResolveUDPAddr("udp", service_address)
	if err != nil {
		return err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, gaddr)
	if err != nil {
		return err
	}

	*res = Responder{gaddr, conn, make(chan cmd), make(map[string]bool)}

	go res.responderLoop()
	go res.socketLoop()

	return nil
}

func (res *Responder) Stop() {
	res.cmdChan <- cmd{op: opQuit}
}

func (res *Responder) Register(service string) {
	res.cmdChan <- cmd{op: opRegister, serv: service}
}

func (res *Responder) Unregister(service string) {
	res.cmdChan <- cmd{op: opUnregister, serv: service}
}

func (res *Responder) responderLoop() {
	defer res.conn.Close()
	defer close(res.cmdChan)

	for {
		cmd := <-res.cmdChan
		switch cmd.op {
		case opQuit:
			return
		case opRequest:
			prefix := "service_request:"
			if strings.HasPrefix(cmd.serv, prefix) {
				cmd.serv = cmd.serv[len(prefix):]
				if res.table[cmd.serv] {
					_, err := res.conn.WriteTo([]byte("service_offer:"+cmd.serv), res.gaddr)
					if err != nil {
						panic(err)
					}
				}
			}
		case opRegister:
			res.table[cmd.serv] = true
		case opUnregister:
			delete(res.table, cmd.serv)
		}
	}
}

func (res *Responder) socketLoop() {
	buf := make([]byte, 2048)
	for {
		n, err := res.conn.Read(buf)
		if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == closed_connection_error {
			return
		} else if err != nil {
			panic(err)
		}
		res.cmdChan <- cmd{opRequest, string(buf[:n])}
	}
}

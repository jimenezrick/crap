package network

import (
	"log"
	"net"
)

import (
	"crap/config"
	"crap/store"
)

type server struct {
	net.Listener
}

func NewServer() *server {
	return new(server)
}

func (s *server) Start() error {
	lis, err := net.Listen("tcp", config.GetString("network.listen_address"))
	if err != nil {
		return err
	}

	go func() {
		for {
			sock, err := lis.Accept()
			if IsErrClosing(err) {
				return
			} else if err != nil {
				panic(err)
			}

			conn := newConn(sock)
			go conn.handleConnection()
		}
	}()

	s.Listener = lis
	return nil
}

func (s *server) Stop() error {
	return s.Close()
}

// XXX XXX XXX
func (c conn) handleConnection() {
	defer c.Close()

	var req request
	if err := c.ReadJSONFrame(&req); err != nil {
		log.Print("Error:", err)
		return
	}
	log.Print("Request:", req)

	switch req.Req {
	case "store":
		key, err := c.handleStore(req)
		if err != nil {
			log.Print("Error:", err)
		}
		log.Print("Key:", key)
	default:
		log.Print("Error: not implemented")
		return
	}
}

func (c conn) handleStore(req request) (string, error) {
	blob, err := store.NewBlob()
	if err != nil {
		return "", err
	}

	err = c.ReadBlobFrameTo(blob)
	if err != nil {
		blob.Abort()
		return "", err
	}

	key, err := blob.Store()
	if err != nil {
		blob.Abort()
		return "", err
	}

	// XXX: Check req.Key with key

	return key, nil
}

// XXX XXX XXX

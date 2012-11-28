package network

// XXX: Add more logging everywhere

import (
	"log" // XXX
	"net"
)

import (
	"crap/kvmap"
	"crap/store"
)

type Network struct {
	address  string
	store    *store.Store
	listener net.Listener
}

func New(config *kvmap.KVMap, store *store.Store) *Network {
	addr, err := config.GetString("network.listen_address")
	if err != nil {
		panic(err)
	}

	return &Network{addr, store, nil}
}

func (n *Network) Start() error {
	lis, err := net.Listen("tcp", n.address)
	if err != nil {
		return err
	}
	n.listener = lis

	go func() {
		for {
			sock, err := n.listener.Accept()
			if IsClosing(err) {
				return
			} else if err != nil {
				panic(err)
			}

			conn := newConn(sock)
			go n.handleConnection(conn)
		}
	}()

	return nil
}

func (n *Network) Stop() error {
	return n.listener.Close()
}

// XXX XXX XXX
func (n *Network) handleConnection(conn *Conn) {
	defer conn.Close()

	var req request
	if err := conn.ReadJSONFrame(&req); err != nil {
		log.Print("Error:", err)
		return
	}
	log.Print("Request:", req)

	switch req.Val {
	case "store":
		key, err := n.handleStore(conn)
		if err != nil {
			log.Print("Error:", err)
		}
		log.Print("Key:", key)
	default:
		log.Print("Error: not implemented")
		return
	}

	if err := conn.ReadJSONFrame(&req); err != nil {
		log.Print("Error:", err)
		return
	}
	log.Print("Request:", req)

	res := result{"ok", "everything went smooth"}
	if err := conn.WriteJSONFrame(res); err != nil {
		log.Print("Error:", err)
		return
	}
}

func (n *Network) handleStore(conn *Conn) (string, error) {
	blob, err := n.store.NewBlob()
	if err != nil {
		return "", err
	}

	if err = conn.ReadBlobFrameTo(blob); err != nil {
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

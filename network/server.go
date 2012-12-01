package network

// XXX: Recover from panic

import (
	"os"
	"net"
)

import (
	"crap/config"
	"crap/store"
	"crap/log"
	"crap/util"
)

func Listen(config config.Config, store *store.Store) (net.Listener, error) {
	lis, err := net.Listen("tcp", config.GetString("network.listen_address"))
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			sock, err := lis.Accept()
			if IsClosing(err) {
				return
			} else if err != nil {
				panic(err)
			}

			conn := newConn(store, sock)
			go conn.handleConnection()
		}
	}()

	return lis, nil
}

func (c *Conn) handleConnection() {
	defer c.Close()
	log.Info.Printf("Connection opened from %s", c.sock.RemoteAddr())

	var req request
	if err := c.ReadJSONFrame(&req); err != nil {
		c.respondError(err)
		panic(err)
	}

	switch req.Val {
	case "store":
		if err := c.handleStore(); os.IsExist(err) {
			c.respondBlobExist()
		} else if err != nil {
			c.respondError(err)
			panic(err)
		}
	default:
		c.respondInvalidRequest()
	}
}

func (c *Conn) handleStore() error {
	log.Info.Printf("Store request from %s", c.sock.RemoteAddr())

	blob, err := c.store.NewBlob()
	if err != nil {
		return err
	}

	abort := true
	defer func() {
		if abort { blob.Abort() }
	}()

	if err := c.ReadBlobFrameTo(blob); err != nil {
		return err
	}

	var key keyRequest
	if err := c.ReadJSONFrame(&key); err != nil {
		return err
	}

	blobKey := util.HexHash(blob)
	if key.Val != blobKey {
		c.respondIncorrectKey()
		return nil
	}

	if _, err := blob.Store(); err != nil {
		return err
	}
	log.Info.Printf("Blob %s stored by %s", blobKey, c.sock.RemoteAddr())

	abort = false
	c.respondOK()
	return nil
}

func (c *Conn) respondOK() {
	c.respond("ok", "")
}

func (c *Conn) respondBlobExist() {
	c.respond("blob_exist", "")
}

func (c *Conn) respondIncorrectKey() {
	c.respond("incorrect_key", "")
}

func (c *Conn) respondInvalidRequest() {
	c.respond("invalid_request", "")
}

func (c *Conn) respondError(err error) {
	c.respond("error", err.Error())
}

func (c *Conn) respond(res, info string) {
	if err := c.WriteJSONFrame(response{res, info}); err != nil {
		panic(err)
	}
}

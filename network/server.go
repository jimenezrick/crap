package network

import (
	"io"
	"bufio"
	"net"
	"log"
)

import (
	"crap/config"
	"crap/store"
)

type Server struct {
	listener net.Listener
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", config.GetString("net.listen_address"))
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := lis.Accept()
			if IsErrClosing(err) {
				return
			} else if err != nil {
				panic(err)
			}

			go handleConnection(conn)
		}
	}()

	s.listener = lis
	return nil
}

func (s *Server) Stop() error {
	return s.listener.Close()
}

// XXX XXX XXX
func handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	var req Request

	if err := ReadJSONFrame(rw, &req); err != nil {
		log.Print("error:", err)
		return
	}

	log.Print("request:", req)

	switch req.Request {
	case "store":
		handleStore(req, rw)
	default:
		log.Print("ERROR")
	}

	conn.Close()
}






func handleStore(req Request, r io.Reader) (string, error) {
	blob, err := store.NewBlob()
	if err != nil {
		return "", err
	}

	err = CopyBlobFrame(blob, r)
	if err != nil {
		blob.Abort()
		return "", err
	}

	key, err := blob.Store()
	if err != nil {
		blob.Abort()
		return "", err
	}

	return key, nil
}
// XXX XXX XXX




type Request struct {
	Request string
	Key string
}

type Result struct {
	Result  string
	Info string
}

//func (r *Request) sanitizeRequest() error {
//        if r.Request != "store" {
//                return errors.New("invalid request")
//        }
//}



// {
//         request: store
//         key: 43874536563475783475374
//         result: ok | error
//         info: No such blob
// }

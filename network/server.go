package network

import (
	"io"
	"bufio"
	"net"
	"log"
)

import "crap/config"

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
	var req request

	bufConn := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	if err := ReadJSONFrame(bufConn, &req); err != nil {
		log.Print("error:", err)
		return
	}

	log.Print("request:", req)

	switch req.Req {
	case "store":
		handleStore(req, bufConn)
	default:
		log.Print("ERROR")
	}

	conn.Close()
}



func handleStore(req request, conn io.ReadWriter) error {
	return nil // XXX
}



// XXX XXX XXX




type request struct {
	Req string
	Key string
}

type result struct {
	Res  string
	Info string
}

//func (r *request) sanitizeRequest() error {
//        if r.Req != "store" {
//                return errors.New("invalid request")
//        }
//}



// {
//         request: store
//         key: 43874536563475783475374
//         result: ok | error
//         info: No such blob
// }

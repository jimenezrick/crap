package main

import (
	"os"
	"time"
)

import (
	"crap/config"
	"crap/network"
	"crap/store"
)

func init() {
	config.Set("store.path", "/tmp")
	config.Set("store.permissions", 0700)

	config.Set("network.listen_address", ":9000")
	config.Set("network.max_json_frame_size", 4096)
}

func main() {
	if len(os.Args) == 2 {
		client()
	} else {
		server()
	}
}

func server() {
	store.Init()

	var serv network.Server

	if err := serv.Start(); err != nil {
		panic(err)
	}

	time.Sleep(1000 * time.Second)
	serv.Stop()
	time.Sleep(3 * time.Second)
}

func client() {
	conn, err := network.Connect("localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	req := network.Request{"store", "foobar"}
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	network.WriteJSONFrame(conn, req)
	network.WriteBlobFrameFrom(conn, file, uint32(info.Size()))
}

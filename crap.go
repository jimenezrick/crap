package main

import "time"

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
	store.Init()

	var serv network.Server

	if err := serv.Start(); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)
	serv.Stop()
	time.Sleep(3 * time.Second)
}

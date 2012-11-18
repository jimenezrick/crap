package main

import "time"

import (
	"crap/config"
	"crap/network"
)

func init() {
	config.Set("network.listen_address",      ":9000")
	config.Set("network.max_json_frame_size", 4096)
}

func main() {
	var serv network.Server

	if err := serv.Start(); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)
	serv.Stop()
	time.Sleep(3 * time.Second)
}

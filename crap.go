package main

import (
	"os"
	"time"
	"fmt"
)

import (
	"crap/kvmap"
	"crap/network"
	"crap/store"
)

func main() {
	if len(os.Args) == 1 {
		server()
	} else if len(os.Args) == 3 {
		client()
	} else {
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "[<host> <file>]")
}

func defaultConfig() *KVMap {
	config := kvmap.New()
	config.Set("store.path", "/tmp")
	config.Set("store.permissions", 0700)
	config.Set("network.listen_address", ":9000")
	config.Set("network.max_json_frame_size", 4096)
}

func loadConfigFile(name string) *KVMap {
	config := defaultConfig()
	configFile, err := kvmap.LoadJSONFile(name)
	if err == nil {
		config.Merge(configFile)
	}

	return config
}


















func server() {
	store, err := store.New()
	if err != nil {
		panic(err)
	}

	serv := network.NewServer()
	if err := serv.Start(); err != nil {
		panic(err)
	}

	time.Sleep(100 * time.Second)
	if err := serv.Stop(); err != nil {
		panic(err)
	}
	time.Sleep(3 * time.Second)
}

func client() {
	conn, err := network.Connect(os.Args[1])
	if err != nil {
		panic(err)
	}
	if err = conn.StoreBlob(os.Args[2]); err != nil {
		panic(err)
	}
	if err = conn.Close(); err != nil {
		panic(err)
	}
}

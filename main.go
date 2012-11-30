package main

import (
	"fmt"
	"os"
	"time"
)

import (
	"crap/network"
	"crap/store"
	"crap/log"
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

func server() {
	config := loadConfig()
	log.Init(config)

	store, err := store.New(config)
	if err != nil {
		panic(err)
	}

	if err := store.Lock(); err != nil {
		panic(err)
	}
	defer store.Unlock()

	network := network.New(config, store)
	if err := network.Start(); err != nil {
		panic(err)
	}
	defer network.Stop()

	time.Sleep(60 * time.Second)
}

func client() {
	config := loadConfig()
	log.Init(config)

	conn, err := network.Connect(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	file, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = conn.StoreBlob(file)
	if err != nil {
		panic(err)
	}
}

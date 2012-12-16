package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

import (
	"crap/config"
	"crap/log"
	"crap/network"
	"crap/store"
)

var configFiles []string = []string{"/etc/crap/config.json", "config.json"}

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
	os.Exit(1)
}

func server() {
	config := config.LoadConfig(configFiles)
	log.Init(config)

	store, err := store.Open(config)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	_, err = network.Listen(config, store)
	if err != nil {
		panic(err)
	}

	handleSignals()
}

func client() {
	config := config.LoadConfig(configFiles)
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

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

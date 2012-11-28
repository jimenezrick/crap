package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

import (
	"crap/kvmap"
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

func defaultConfig() *kvmap.KVMap {
	config := kvmap.New()
	config.Set("log.debug", true)
	config.Set("log.syslog", true)
	config.Set("store.path", "/tmp")
	config.Set("store.permissions", 0700)
	config.Set("network.listen_address", ":9000")
	return config
}

func loadConfigFile(name string) *kvmap.KVMap {
	config := defaultConfig()
	configFile, err := kvmap.LoadJSONFile(name)
	if err == nil {
		config.Merge(configFile)
	}
	return config
}

func server() {
	config := loadConfigFile("crap.conf")
	log.Init(config)

	store, err := store.New(config)
	if err != nil {
		panic(err)
	}

	network := network.New(config, store)
	if err := network.Start(); err != nil {
		panic(err)
	}

	time.Sleep(100 * time.Second)
	if err := network.Stop(); err != nil {
		panic(err)
	}
	time.Sleep(3 * time.Second)
}

func client() {
	config := loadConfigFile("crap.conf")
	log.Init(config)

	conn, err := network.Connect(os.Args[1])
	if err != nil {
		panic(err)
	}

	blob, size, err := takeBlob(os.Args[2])
	if err != nil {
		panic(err)
	}

	if err = conn.StoreBlob(blob, uint32(size)); err != nil {
		panic(err)
	}

	if err = conn.Close(); err != nil {
		panic(err)
	}
}

// XXX XXX XXX
func takeBlob(name string) (io.Reader, int64, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, 0, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	return bufio.NewReader(file), info.Size(), nil
}

// XXX XXX XXX

package main

import (
	"os"
	"fmt"
	"math/rand"
)

import "crap/skiplist"

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	var size int
	if _, err := fmt.Sscan(os.Args[1], &size); err != nil {
		panic(err)
	}

	s := skiplist.NewIntMap()
	for i := 0; i < size; i++ {
		s.Insert(rand.Intn(size), i)
	}
	s.DumpASCII(os.Stderr)
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "<size>")
	os.Exit(1)
}

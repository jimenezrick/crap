package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
)

import "crap/skiplist"

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	var size uint
	if _, err := fmt.Sscan(os.Args[2], &size); err != nil {
		panic(err)
	}

	if os.Args[1] == "-s" {
		testSkipList(size)
	} else if os.Args[1] == "-m" {
		testMap(size)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "-s|-m <size>")
	os.Exit(1)
}

func testSkipList(size uint) {
	s := skiplist.NewIntMap()
	s.SetMaxLevel(size)

	for ; size > 0; size-- {
		s.Insert(rand.Int(), nil)
	}

	fmt.Println("Skip List with length", s.Len())
	printMemStats()
}

func testMap(size uint) {
	m := make(map[int]interface{}, size)

	for ; size > 0; size-- {
		m[rand.Int()] = nil
	}

	fmt.Println("Map with length", len(m))
	printMemStats()
}

func printMemStats() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	fmt.Println("Bytes allocated", stats.Alloc)

	fmt.Println("Press ENTER")
	os.Stdin.Read([]byte{0})
}

package main

import (
	"os"
	"crap/skiplist"
)

func main() {
	s := skiplist.NewIntMap()
	s.Insert(2, "foo")
	s.Insert(1, "bar")
	s.DumpASCII(os.Stderr)
}

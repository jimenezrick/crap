package main

import (
	"os"
	"crap/skiplist"
)

func main() {
	s := skiplist.New()
	s.DumpASCII(os.Stderr)
}

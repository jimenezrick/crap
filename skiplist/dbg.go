package skiplist

import (
	"fmt"
	"io"
)

func (s *SkipList) DumpASCII(w io.Writer) {
	fmt.Fprintln(w, "--- Skip List ---")
	fmt.Fprintln(w, "length   = %d", s.length)
	fmt.Fprintln(w, "p        = %f", s.p)
	fmt.Fprintln(w, "maxLevel = %d", s.maxLevel)
	fmt.Fprintln(w, "-----------------")
	fmt.Fprintln(w)

	node := s.header
	levels := len(node.forward)

	var length uint
	for length = 0; node != nil; length, node = length + 1, node.forward[0] {
		printNode(w, node, levels)
		printLinks(w, levels)

		if len(node.forward) > levels {
			panic("skiplist: levels mismatch")
		}
	}

	if length - 1 != s.length {
		panic("skiplist: length mismatch")
	}
	printNils(w, levels)
}

func printNode(w io.Writer, node *node, levels int) {
	fmt.Fprintf(w, "%v:\t", node.key)
	for _, f := range node.forward {
		if f != nil {
			fmt.Fprintf(w, "[%v]\t", f.value)
		} else {
			fmt.Fprintf(w, "[*]\t")
		}
	}
	printLinks(w, levels - len(node.forward))
}

func printLinks(w io.Writer, levels int) {
	fmt.Fprint(w, "\t")
	for i := 0; i < levels; i++ {
		fmt.Fprint(w, " |\t")
	}
	fmt.Fprintln(w)
}

func printNils(w io.Writer, levels int) {
	fmt.Fprint(w, "\t")
	for i := 0; i < levels; i++ {
		fmt.Fprint(w, "nil\t")
	}
	fmt.Fprintln(w)
}

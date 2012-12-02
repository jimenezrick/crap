package skiplist

import "math/rand"

type node struct {
	forward    []*node
	key, value interface{}
}

func (n *node) next() *node {
	if len(n.forward) == 0 {
		return nil
	}
	return n.forward[0]
}

type SkipList struct {
	less     func(l, r interface{}) bool
	header   *node
	length   uint
	p        float64
	maxLevel int
	rand     *rand.Rand
}

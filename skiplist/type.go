package skiplist

import (
	"bytes"
	"math/rand"
	"time"
)

type Key interface{}

type Value interface{}

type LessFunc func(l, r Key) bool

type node struct {
	forward []*node
	key     Key
	value   Value
}

func (n *node) next() *node {
	return n.forward[0]
}

type SkipList struct {
	less     LessFunc
	header   *node
	length   uint
	p        float64
	maxLevel int
	rand     *rand.Rand
}

type Ordered interface {
	Less(Ordered) bool
}

func NewCustomMap(less LessFunc) *SkipList {
	return &SkipList{
		less:     less,
		header:   &node{forward: []*node{nil}},
		p:        p,
		maxLevel: maxLevel,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func New() *SkipList {
	return NewCustomMap(func(l, r Key) bool {
		return l.(Ordered).Less(r.(Ordered))
	})
}

func NewIntMap() *SkipList {
	return NewCustomMap(func(l, r Key) bool {
		return l.(int) < r.(int)
	})
}

func NewStringMap() *SkipList {
	return NewCustomMap(func(l, r Key) bool {
		return l.(string) < r.(string)
	})
}

func NewByteMap() *SkipList {
	return NewCustomMap(func(l, r Key) bool {
		return bytes.Compare(l.([]byte), r.([]byte)) == -1
	})
}

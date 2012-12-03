package skiplist

import (
	"bytes"
	"math/rand"
	"time"
)

type node struct {
	forward    []*node
	key, value interface{}
}

func (n *node) next() *node {
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

type Ordered interface {
	Less(Ordered) bool
}

func NewCustomMap(less func(l, r interface{}) bool) *SkipList {
	return &SkipList{
		less:     less,
		header:   &node{forward: []*node{nil}},
		p:        p,
		maxLevel: maxLevel,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func New() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(Ordered).Less(r.(Ordered))
	})
}

func NewIntMap() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(int) < r.(int)
	})
}

func NewStringMap() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return l.(string) < r.(string)
	})
}

func NewByteMap() *SkipList {
	return NewCustomMap(func(l, r interface{}) bool {
		return bytes.Compare(l.([]byte), r.([]byte)) == -1
	})
}

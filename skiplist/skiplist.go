package skiplist

// XXX: Benchmark fixed dice
// XXX: Dump dot file

// Ref: William Pugh "Skip lists: a probabilistic alternative to balanced trees"

import "math/rand"

func (s *SkipList) SetP(p float64) {
	if s.length != 0 {
		panic("skiplist: container not empty")
	}
	s.p = p
}

func (s *SkipList) SetMaxLevel(maxElems uint) {
	if s.length != 0 {
		panic("skiplist: container not empty")
	}
	s.maxLevel = expectedLevels(s.p, maxElems)
}

func (s *SkipList) SetRand(src rand.Source) {
	if s.length != 0 {
		panic("skiplist: container not empty")
	}
	s.rand = rand.New(src)
}

func (s *SkipList) Len() uint {
	return s.length
}

func (s *SkipList) level() int {
	return len(s.header.forward) - 1
}

func (s *SkipList) randomLevel() (n int) {
	level := s.level()
	for n = 0; s.rand.Float64() < s.p && n <= level && n < s.maxLevel; n++ {
	}
	return
}

func (s *SkipList) getPath(update []*node, key interface{}) *node {
	//
	// XXX
	//

	return nil
}

// XXX XXX XXX

/*
func New() {
	src := rand.NewSource(time.Now().UnixNano())

	src.Uint32()
}
*/

//func SetSeed()

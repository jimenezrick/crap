// Read "Skip lists: a probabilistic alternative to balanced trees" and
// "A skip list cookbook" by William Pugh.
package skiplist

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
	// Returns the level-1 of the skip list, used for slices indices.
	// The level of an empty skip list is 1.
	return len(s.header.forward) - 1
}

func (s *SkipList) randomLevel() (n int) {
	// Returns a random level in the range [0, s.level()+1] been at most
	// equal to s.maxLevel-1. Used for slices indices.
	for n = 0; s.rand.Float64() < s.p && n < s.maxLevel-1; n++ {
	}
	return
}

func (s *SkipList) getPath(update []*node, key Key, next bool) *node {
	current := s.header
	for i := s.level(); i >= 0; i-- {
		for current.forward[i] != nil && s.less(current.forward[i].key, key) {
			current = current.forward[i]
		}

		if update != nil {
			update[i] = current
		}
	}

	if next {
		return current.next()
	} else if current == s.header {
		return nil
	}
	return current
}

func (s *SkipList) Lookup(key Key) (Value, bool) {
	candidate := s.getPath(nil, key, true)
	if candidate == nil || candidate.key != key {
		return nil, false
	}
	return candidate.value, true
}

func (s *SkipList) Lesser(max Key) (Key, Value) {
	candidate := s.getPath(nil, max, false)
	if candidate != nil {
		return candidate.key, candidate.value
	}
	return nil, nil
}

func (s *SkipList) GreaterOrEqual(min Key) (Key, Value) {
	candidate := s.getPath(nil, min, true)
	if candidate != nil {
		return candidate.key, candidate.value
	}
	return nil, nil
}

func (s *SkipList) Min() (Key, Value) {
	min := s.header.next()
	if min != nil {
		return min.key, min.value
	}
	return nil, nil
}

func (s *SkipList) Max() (Key, Value) {
	current := s.header
	for i := s.level(); i >= 0; i-- {
		for current.forward[i] != nil {
			current = current.forward[i]
		}
	}

	if current == s.header {
		return nil, nil
	}
	return current.key, current.value
}

func (s *SkipList) Insert(key Key, value Value) (Value, bool) {
	return s.insert(key, value, false)
}

func (s *SkipList) InsertMulti(key Key, value Value) (Value, bool) {
	return s.insert(key, value, true)
}

func (s *SkipList) insert(key Key, value Value, multi bool) (Value, bool) {
	update := make([]*node, s.level()+1)
	candidate := s.getPath(update, key, true)
	if candidate != nil && candidate.key == key && !multi {
		oldValue := candidate.value
		candidate.value = value
		return oldValue, true
	}

	newLevel := s.randomLevel()
	if level := s.level(); newLevel > level {
		for i := level + 1; i <= newLevel; i++ {
			update = append(update, s.header)
			s.header.forward = append(s.header.forward, nil)
		}
	}

	node := &node{make([]*node, newLevel+1), key, value}
	for i := 0; i <= newLevel; i++ {
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node
	}

	s.length++
	return nil, false
}

func (s *SkipList) Delete(key Key) (Value, bool) {
	update := make([]*node, s.level()+1)
	candidate := s.getPath(update, key, true)
	if candidate == nil || candidate.key != key {
		return nil, false
	}

	for i, level := 0, s.level(); i <= level &&
		update[i].forward[i] == candidate; i++ {
		update[i].forward[i] = candidate.forward[i]
	}

	for s.level() > 0 && s.header.forward[s.level()] == nil {
		s.header.forward = s.header.forward[:s.level()]
	}

	s.length--
	return candidate.value, true
}

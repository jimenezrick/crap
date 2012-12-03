package skiplist

type Set struct {
	SkipList
}

func NewSet() *Set {
	return &Set{*New()}
}

func NewIntSet() *Set {
	return &Set{*NewIntMap()}
}

func NewStringSet() *Set {
	return &Set{*NewStringMap()}
}

func NewByteSet() *Set {
	return &Set{*NewByteMap()}
}

func (s *Set) Insert(key Key) bool {
	_, ok := s.SkipList.Insert(key, nil)
	return ok
}

func (s *Set) InsertMulti(key Key) bool {
	_, ok := s.SkipList.Insert(key, nil)
	return ok
}

func (s *Set) Contains(key Key) bool {
	_, ok := s.SkipList.Get(key)
	return ok
}
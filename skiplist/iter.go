package skiplist

type Iterator interface {
	Next() bool
	Key() interface{}
	Value() interface{}
}

type iterator struct {
	current *node
	key, value interface{}
}

func (i *iterator) Next() bool {
	if next := i.current.next(); next != nil {
		i.current = next
		i.key = next.key
		i.value = next.value
		return true
	}
	return false
}

func (i *iterator) Key() interface{} {
	return i.key
}

func (i *iterator) Value() interface{} {
	return i.value
}

type rangeIterator struct {
	iterator
	limit interface{}
	skipList *SkipList
}

func (r *rangeIterator) Next() bool {
	if next := r.current.next(); next != nil {
		if r.skipList.less(next.key, r.limit) {
			r.current = next
			r.key = next.key
			r.value = next.value
			return true
		}
	}
	return false
}

func (s *SkipList) Iterator() Iterator {
	return &iterator{current: s.header}
}

func (s *SkipList) Range(from, to interface{}) Iterator {
	start := s.getPath(nil, from)
	return &rangeIterator{
		iterator: iterator{current: &node{forward: []*node{start}}},
		limit: to,
		skipList: s,
	}
}

package skiplist

type Iterator interface {
	Next() (key Key, value Value)
}

type iterator struct {
	current *node
}

func (i *iterator) Next() (Key, Value) {
	if next := i.current.next(); next != nil {
		i.current = next
		return next.key, next.value
	}
	return nil, nil
}

type rangeIterator struct {
	iterator
	limit    Key
	skipList *SkipList
}

func (r *rangeIterator) Next() (Key, Value) {
	if next := r.current.next(); next != nil {
		if r.skipList.less(next.key, r.limit) {
			r.current = next
			return next.key, next.value
		}
	}
	return nil, nil
}

func (s *SkipList) Iterator() Iterator {
	return &iterator{s.header}
}

func (s *SkipList) Range(from, to Key) Iterator {
	start := s.getPath(nil, from, true)
	return &rangeIterator{
		iterator: iterator{&node{forward: []*node{start}}},
		limit:    to,
		skipList: s,
	}
}

package skiplist

type Iterator interface {
	Next() (key interface{}, value interface{})
}

type iterator struct {
	current    *node
	key, value interface{}
}

func (i *iterator) Next() (interface{}, interface{}) {
	if next := i.current.next(); next != nil {
		i.current = next
		return next.key, next.value
	}
	return nil, nil
}

type rangeIterator struct {
	iterator
	limit    interface{}
	skipList *SkipList
}

func (r *rangeIterator) Next() (interface{}, interface{}) {
	if next := r.current.next(); next != nil {
		if r.skipList.less(next.key, r.limit) {
			r.current = next
			return next.key, next.value
		}
	}
	return nil, nil
}

func (s *SkipList) Iterator() Iterator {
	return &iterator{current: s.header}
}

func (s *SkipList) Range(from, to interface{}) Iterator {
	start := s.getPath(nil, from, true)
	return &rangeIterator{
		iterator: iterator{current: &node{forward: []*node{start}}},
		limit:    to,
		skipList: s,
	}
}

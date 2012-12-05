package skiplist

type Heap struct {
	SkipList
}

func NewHeap() *Heap {
	return &Heap{*New()}
}

func NewIntHeap() *Heap {
	return &Heap{*NewIntMap()}
}

func (h *Heap) Pop() Key {
	first := h.header.next()
	if first == nil {
		return nil
	}

	for i, level := 0, h.level(); i <= level && h.header.forward[i] == first; i++ {
		h.header.forward[i] = first.forward[i]
	}
	return first.key
}

func (h *Heap) Peek() Key {
	key, _ := h.SkipList.Min()
	return key
}

func (h *Heap) Push(item Key) {
	h.SkipList.InsertMulti(item, nil)
}

func (h *Heap) Contains(item Key) bool {
	_, ok := h.SkipList.Lookup(item)
	return ok
}

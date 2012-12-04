package skiplist

import (
	"math/rand"
	"testing"
)

func TestSetMaxLevel(t *testing.T) {
	s := New()
	s.SetP(0.25)
	s.SetMaxLevel(4294967296 - 1)
	if s.maxLevel != 16 {
		t.Error("maxLevel should be 16")
	}

	s.SetP(0.5)
	s.SetMaxLevel(4294967296 - 1)
	if s.maxLevel != 32 {
		t.Error("maxLevel should be 32")
	}
}

func TestLen(t *testing.T) {
	s := NewIntMap()
	if s.Len() != 0 {
		t.Error("length should be 0")
	}

	s.Insert(1, 2)
	s.Insert(1, 123)
	s.InsertMulti(1, 666)
	if s.Len() != 2 {
		t.Error("length should be 2")
	}

	if val, ok := s.Delete(1); !ok || val != 666 {
		t.Error("value should be 666")
	}
	if val, ok := s.Delete(1); !ok || val != 123 {
		t.Error("value should be 123")
	}
	if val, ok := s.Delete(1); ok || val != nil {
		t.Error("value should be nil")
	}
}

func TestLevel(t *testing.T) {
	s := New()
	if s.level() != 0 {
		t.Error("level should be 0")
	}
}

/*
func TestRandomLevel(t *testing.T) {
	s := New()
	s.SetP(1.0)
	s.maxLevel = 3
	if s.randomLevel() != 2 {
		t.Error("random level should be 2")
	}
}
*/

func TestGetPath(t *testing.T) {
	s := NewIntMap()
	if s.getPath(nil, nil, false) != nil {
		t.Error("node should be nil")
	}
	if s.getPath(nil, nil, true) != nil {
		t.Error("node should be nil")
	}

	s.Insert(1, nil)
	update := make([]*node, s.level()+1)

	if s.getPath(update, 1, false) != nil {
		t.Error("node should be nil")
	}
	if update[0] != s.header {
		t.Error("update should contain header")
	}

	if s.getPath(update, 1, true) != s.header.next() {
		t.Error("node should be s.header.next()")
	}
	if update[0] != s.header {
		t.Error("update should contain header")
	}

	if s.getPath(update, 2, true) != nil {
		t.Error("node should be nil")
	}
	if update[0] != s.header.next() {
		t.Error("update should contain header")
	}

	if s.getPath(update, 2, false) != s.header.next() {
		t.Error("node should be s.header.next()")
	}
	if update[0] != s.header.next() {
		t.Error("update should contain header")
	}
}

func TestGet(t *testing.T) {
	s := NewIntMap()
	if val, ok := s.Get(1); ok || val != nil {
		t.Error("value should be nil")
	}

	if val, ok := s.Insert(1, 123); ok || val != nil {
		t.Error("value should be nil")
	}
	if val, ok := s.Insert(1, 666); !ok || val != 123 {
		t.Error("value should be 123")
	}

	if val, ok := s.InsertMulti(2, 123); ok || val != nil {
		t.Error("value should be nil")
	}
	if val, ok := s.InsertMulti(2, 666); ok || val != nil {
		t.Error("value should be nil")
	}

	if val, ok := s.Get(1); !ok || val != 666 {
		t.Error("value should be 666")
	}
	if val, ok := s.Get(2); !ok || val != 666 {
		t.Error("value should be 666")
	}
}

func TestGetLesser(t *testing.T) {
	s := NewIntMap()
	if key, val := s.GetLesser(1); key != nil || val != nil {
		t.Error("lesser should be nil")
	}

	s.Insert(1, nil)
	if key, val := s.GetLesser(1); key != nil || val != nil {
		t.Error("lesser should be nil")
	}

	s.Insert(2, nil)
	if key, val := s.GetLesser(2); key != 1 || val != nil {
		t.Error("lesser should be 1")
	}

	if key, val := s.GetLesser(8); key != 2 || val != nil {
		t.Error("lesser should be 2")
	}
}

func TestGetGreaterOrEqual(t *testing.T) {
	s := NewIntMap()
	if key, val := s.GetGreaterOrEqual(1); key != nil || val != nil {
		t.Error("greater or equal should be nil")
	}

	s.Insert(1, nil)
	if key, val := s.GetGreaterOrEqual(8); key != nil || val != nil {
		t.Error("greater or equal should be nil")
	}

	s.Insert(8, nil)
	if key, val := s.GetGreaterOrEqual(8); key != 8 || val != nil {
		t.Error("greater or equal should be 8")
	}

	if key, val := s.GetGreaterOrEqual(7); key != 8 || val != nil {
		t.Error("greater or equal should be 8")
	}
}

func TestGetMin(t *testing.T) {
	s := NewIntMap()
	if key, val := s.GetMin(); key != nil || val != nil {
		t.Error("min should be nil")
	}

	s.Insert(3, nil)
	s.Insert(8, nil)
	s.Insert(2, nil)
	if key, val := s.GetMin(); key != 2 || val != nil {
		t.Error("min should be 2")
	}
}

func TestGetMax(t *testing.T) {
	s := NewIntMap()
	if key, val := s.GetMax(); key != nil || val != nil {
		t.Error("max should be nil")
	}

	s.Insert(3, nil)
	s.Insert(8, nil)
	s.Insert(2, nil)
	if key, val := s.GetMax(); key != 8 || val != nil {
		t.Error("max should be 8")
	}
}

func TestDelete(t *testing.T) {
	s := NewIntMap()
	s.Insert(3, 1)
	s.Insert(8, 2)
	s.InsertMulti(8, 4)

	if val, ok := s.Delete(100); ok || val != nil {
		t.Error("value should be nil")
	}

	if val, ok := s.Delete(3); !ok || val != 1 {
		t.Error("value should be 1")
	}
	if val, ok := s.Delete(3); ok || val != nil {
		t.Error("value should be nil")
	}

	if val, ok := s.Delete(8); !ok || val != 4 {
		t.Error("value should be 4")
	}
	if val, ok := s.Delete(8); !ok || val != 2 {
		t.Error("value should be 2")
	}
	if val, ok := s.Delete(8); ok || val != nil {
		t.Error("value should be nil")
	}
}

func TestIterator(t *testing.T) {
	s := NewIntMap()
	s.Insert(1, 123)
	s.Insert(2, 666)
	s.Insert(3, 8)

	iter := s.Iterator()
	for i, v := range []int{123, 666, 8} {
		key, val := iter.Next()
		if key != i+1 || val != v {
			t.Error("value should be", v)
		}
	}
}

func TestRange(t *testing.T) {
	s := NewIntMap()
	s.Insert(1, 123)
	s.Insert(2, 666)
	s.Insert(3, 8)
	s.Insert(4, 999)
	s.Insert(5, 0)

	iter := s.Range(2, 4)
	for i, v := range []int{666, 8} {
		key, val := iter.Next()
		if key != i+2 || val != v {
			t.Error("value should be", v)
		}
	}
}

func TestHeap(t *testing.T) {
	h := NewIntHeap()
	if h.Peek() != nil {
		t.Error("peek should be nil")
	}

	h.Push(2)
	h.Push(3)
	h.Push(2)
	h.Push(1)

	if h.Pop() != 1 {
		t.Error("pop should be 1")
	}
	if h.Pop() != 2 {
		t.Error("pop should be 2")
	}
	if h.Pop() != 2 {
		t.Error("pop should be 2")
	}
	if h.Pop() != 3 {
		t.Error("pop should be 3")
	}

	if h.Peek() != nil {
		t.Error("peek should be nil")
	}
}

func buildFromList(s *SkipList, ints []int) {
	for _, i := range ints {
		s.Insert(i, nil)
	}
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	ints := rand.Perm(b.N)
	s := NewIntMap()

	b.StartTimer()
	buildFromList(s, ints)
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	ints := rand.Perm(b.N)
	s := NewIntMap()
	buildFromList(s, ints)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Delete(ints[i])
	}
}

func BenchmarkLookup(b *testing.B) {
	b.StopTimer()
	ints := rand.Perm(b.N)
	s := NewIntMap()
	buildFromList(s, ints)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Get(ints[i])
	}
}

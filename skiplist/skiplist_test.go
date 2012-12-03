package skiplist

import "testing"

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

	s.Insert(1, nil)
	s.Insert(1, nil)
	s.InsertMulti(1, nil)
	if s.Len() != 2 {
		t.Error("length should be 2")
	}
	// TODO: Use Delete()
}

func TestLevel(t *testing.T) {
	s := New()
	if s.level() != 0 {
		t.Error("level should be 0")
	}
}

func TestRandomLevel(t *testing.T) {
	s := New()
	s.SetP(1.0)
	s.maxLevel = 3
	if s.randomLevel() != 1 {
		t.Error("random level should be 1")
	}

	s.header.forward = []*node{nil, nil}
	if s.randomLevel() != 2 {
		t.Error("random level should be 2")
	}

	s.header.forward = []*node{nil, nil, nil}
	if s.randomLevel() != 2 {
		t.Error("random level should be 2")
	}
}

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

	if s.getPath(update, 1, true) != s.header.forward[0] {
		t.Error("node should be s.header.forward[0]")
	}
	if update[0] != s.header {
		t.Error("update should contain header")
	}

	if s.getPath(update, 2, true) != nil {
		t.Error("node should be nil")
	}
	if update[0] != s.header.forward[0] {
		t.Error("update should contain header")
	}

	if s.getPath(update, 2, false) != s.header.forward[0] {
		t.Error("node should be s.header.forward[0]")
	}
	if update[0] != s.header.forward[0] {
		t.Error("update should contain header")
	}
}

func TestGet(t *testing.T) {
	s := NewIntMap()
	if val, ok := s.Get(1); ok || val != nil {
		t.Error("value should be nil")
	}

	if s.Insert(1, 123) != false {
		t.Error("ok should be false")
	}
	if s.Insert(1, 666) != true {
		t.Error("ok should be true")
	}

	if s.InsertMulti(2, 123) != false {
		t.Error("ok should be false")
	}
	if s.InsertMulti(2, 666) != false {
		t.Error("ok should be false")
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

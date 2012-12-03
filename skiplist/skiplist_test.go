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
	s := New()
	if s.Len() != 0 {
		t.Error("length should be 0")
	}
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
	s := New()
	if s.getPath(nil, nil, false) != nil {
		t.Error("node should be nil")
	}

	if s.getPath(nil, nil, true) != nil {
		t.Error("node should be nil")
	}
}

func TestGet(t *testing.T) {
	s := New()
	if val, ok := s.Get(123); ok || val != nil {
		t.Error("value should be nil")
	}
}

func TestGetMin(t *testing.T) {
	s := New()
	if key, val := s.GetMin(); key != nil || val != nil {
		t.Error("min should be nil")
	}
}

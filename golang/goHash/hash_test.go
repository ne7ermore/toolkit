package goHash

import "testing"

func TestHashHas(t *testing.T) {
	hash := New()
	if hash.Has("test") {
		t.Fail()
	}
}

func TestHashAdd(t *testing.T) {
	hash := New()
	hash.Add("test")
	if !hash.Has("test") {
		t.Fail()
	}
}

func TestHashRemove(t *testing.T) {
	hash := New()
	hash.Add("test")
	hash.Remove("test")
	if hash.Has("test") {
		t.Fail()
	}
}

func TestHashAddSameV(t *testing.T) {
	hash := New()
	hash.Add("test")
	hash.Add("test")
	if len(hash.Tuh) != 1 {
		t.Fail()
	}
}

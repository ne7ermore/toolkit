package goHash

import "testing"

func TestHashHasI(t *testing.T) {
	hash := NewI()
	if hash.HasI(123) {
		t.Fail()
	}
}

func TestHashAddI(t *testing.T) {
	hash := NewI()
	hash.AddI(123)
	if !hash.HasI(123) {
		t.Fail()
	}
}

func TestHashRemoveI(t *testing.T) {
	hash := NewI()
	hash.AddI(123)
	hash.RemoveI(123)
	if hash.HasI(123) {
		t.Fail()
	}
}

func TestHashAddSameVI(t *testing.T) {
	hash := NewI()
	hash.AddI(123)
	hash.AddI(123)
	if hash.SizeI() == 2 {
		t.Fail()
	}
}

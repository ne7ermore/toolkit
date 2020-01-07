package goHash

import "testing"
import "fmt"

func TestMapIHas(t *testing.T) {
	m := NewMapI()
	if m.HasI(11111111) {
		t.Fail()
	}
}

func TestMapIAdd(t *testing.T) {
	m := NewMapI()
	m.AddI(222222222, []string{"nihao", "zaijian"})
	if !m.HasI(222222222) {
		t.Fail()
	}
}

func TestMapIRemove(t *testing.T) {
	m := NewMapI()
	vmap := map[string]string{
		"1": "one",
		"2": "two",
	}
	m.AddI(3333333333, vmap)
	if _, ok := m.GetI(3333333333); !ok {
		t.Fail()
	}
	m.RemoveI(3333333333)
	if m.HasI(3333333333) {
		t.Fail()
	}
}

func TestMapIRange(t *testing.T) {
	m := NewMapI()
	m.AddI(1, 1)
	m.AddI(2, 2)
	m.AddI(3, 3)
	for _, i := range m.IterI() {
		fmt.Println(i.Key, i.Val)
	}
	fmt.Println(m.SizeI())
	m.RemoveI(1)
	for _, i := range m.IterI() {
		fmt.Println(i.Key, i.Val)
	}
	fmt.Println(m.SizeI())
}

package goHash

import "testing"
import "fmt"

func TestMapHas(t *testing.T) {
	m := NewMap()
	if m.Has("test") {
		t.Fail()
	}
}

func TestMapAdd(t *testing.T) {
	m := NewMap()
	m.Add("test", []string{"nihao", "zaijian"})
	if !m.Has("test") {
		t.Fail()
	}
}

func TestMapRemove(t *testing.T) {
	m := NewMap()
	vmap := map[string]string{
		"1": "one",
		"2": "two",
	}
	m.Add("test", vmap)
	if _, ok := m.Get("test"); !ok {
		t.Fail()
	}
	m.Remove("test")
	if m.Has("test") {
		t.Fail()
	}
}

func TestMapRange(t *testing.T) {
	m := NewMap()
	m.Add("1", 1)
	m.Add("2", 2)
	m.Add("3", 3)
	for _, i := range m.Iter() {
		fmt.Println(i.Key, i.Val)
	}
}

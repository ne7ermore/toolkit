// Copyright 2017 All rights reserved.
// Author: ne7ermore.

package goHash

import "sync"

type GoMapI struct {
	tuh threadUnsafeHashI
	sync.RWMutex
}

func NewMapI() GoMapI {
	return GoMapI{tuh: newThreadUnsafeHashI()}
}

func (m *GoMapI) AddI(key uint64, value interface{}) bool {
	m.Lock()
	defer m.Unlock()
	return m.tuh.addMapI(key, value)
}

func (m *GoMapI) GetI(key uint64) (interface{}, bool) {
	m.Lock()
	defer m.Unlock()
	return m.tuh.getI(key)
}

func (m *GoMapI) HasI(key uint64) bool {
	m.RLock()
	defer m.RUnlock()
	return m.tuh.hasI(key)
}

func (m *GoMapI) RemoveI(key uint64) {
	m.Lock()
	defer m.Unlock()
	delete(m.tuh, key)
}

func (m *GoMapI) SizeI() int {
	return int(len(m.tuh))
}

type ItemI struct {
	Key uint64
	Val interface{}
}

func (m *GoMapI) IterI() []*ItemI {
	items := make([]*ItemI, 0)
	m.RLock()
	defer m.RUnlock()
	for key, val := range m.tuh {
		items = append(items, &ItemI{key, val})
	}
	return items
}

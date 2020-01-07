// Copyright 2017 All rights reserved.
// Author: ne7ermore.

package goHash

import "sync"

type GoMap struct {
	tuh threadUnsafeHash
	sync.RWMutex
}

func NewMap() GoMap {
	return GoMap{tuh: newThreadUnsafeHash()}
}

func (m *GoMap) Add(key string, value interface{}) bool {
	m.Lock()
	defer m.Unlock()
	return m.tuh.addMap(key, value)
}

func (m *GoMap) Get(key string) (interface{}, bool) {
	m.Lock()
	defer m.Unlock()
	return m.tuh.get(key)
}

func (m *GoMap) Has(key string) bool {
	m.RLock()
	defer m.RUnlock()
	return m.tuh.has(key)
}

func (m *GoMap) Remove(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.tuh, key)
}

func (m *GoMap) Size() int {
	return int(len(m.tuh))
}

type Item struct {
	Key string
	Val interface{}
}

func (m *GoMap) Iter() []*Item {
	items := make([]*Item, 0)
	m.RLock()
	defer m.RUnlock()
	for key, val := range m.tuh {
		items = append(items, &Item{key, val})
	}
	return items
}

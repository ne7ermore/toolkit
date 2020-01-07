package goHash

import "sync"

type GoHashI struct {
	tuh threadUnsafeHashI
	sync.RWMutex
}

func NewI() *GoHashI {
	return &GoHashI{tuh: newThreadUnsafeHashI()}
}

func (hash *GoHashI) AddI(key uint64) bool {
	hash.Lock()
	hash.Unlock()
	return hash.tuh.addHashI(key)
}

func (hash *GoHashI) HasI(key uint64) bool {
	hash.RLock()
	defer hash.RUnlock()
	return hash.tuh.hasI(key)
}

func (hash *GoHashI) RemoveI(key uint64) {
	hash.Lock()
	delete(hash.tuh, key)
	hash.Unlock()
}

func (hash *GoHashI) SizeI() int {
	return int(len(hash.tuh))
}

package goHash

import "sync"

type GoHash struct {
	Tuh threadUnsafeHash
	sync.RWMutex
}

func newHash() GoHash {
	return GoHash{Tuh: newThreadUnsafeHash()}
}

func New() *GoHash {
	return &GoHash{Tuh: newThreadUnsafeHash()}
}

func (hash *GoHash) Add(s string) bool {
	hash.Lock()
	ret := hash.Tuh.addHash(s)
	hash.Unlock()
	return ret
}

func (hash *GoHash) Has(s string) bool {
	hash.RLock()
	defer hash.RUnlock()
	return hash.Tuh.has(s)
}

func (hash *GoHash) Remove(s string) {
	hash.Lock()
	delete(hash.Tuh, s)
	hash.Unlock()
}

func (hash *GoHash) Size() int {
	return int(len(hash.Tuh))
}

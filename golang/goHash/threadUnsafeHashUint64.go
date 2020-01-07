// Copyright 2017 All rights reserved.
// Author: ne7ermore.

package goHash

// Unsafe thread map
// can modify with more than one thread
type threadUnsafeHashI map[uint64]interface{}

func newThreadUnsafeHashI() threadUnsafeHashI {
	return make(threadUnsafeHashI)
}

// return true if s added
// return false if s exited
func (hash *threadUnsafeHashI) addHashI(key uint64) bool {
	_, found := (*hash)[key]
	if !found {
		(*hash)[key] = struct{}{}
	}
	return !found
}

// return true if value added
// return false if key exited
func (hash *threadUnsafeHashI) addMapI(key uint64, value interface{}) bool {
	_, found := (*hash)[key]
	if !found {
		(*hash)[key] = value
	}
	return !found
}

//Return true if it existed already
func (hash *threadUnsafeHashI) hasI(key uint64) bool {
	_, found := (*hash)[key]
	return found
}

// return val and true or nil and false
func (hash *threadUnsafeHashI) getI(key uint64) (interface{}, bool) {
	value, found := (*hash)[key]
	return value, found
}

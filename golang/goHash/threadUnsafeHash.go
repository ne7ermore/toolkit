// Copyright 2017 All rights reserved.
// Author: ne7ermore.

package goHash

// Unsafe thread map
// can modify with more than one thread
type threadUnsafeHash map[string]interface{}

func newThreadUnsafeHash() threadUnsafeHash {
	return make(threadUnsafeHash)
}

// return true if s added
// return false if s exited
func (hash *threadUnsafeHash) addHash(s string) bool {
	_, found := (*hash)[s]
	if !found {
		(*hash)[s] = struct{}{}
	}
	return !found
}

// return true if value added
// return false if key exited
func (hash *threadUnsafeHash) addMap(key string, value interface{}) bool {
	_, found := (*hash)[key]
	if !found {
		(*hash)[key] = value
	}
	return !found
}

//Return true if it existed already
func (hash *threadUnsafeHash) has(s string) bool {
	_, found := (*hash)[s]
	return found
}

// return val and true or nil and false
func (hash *threadUnsafeHash) get(s string) (interface{}, bool) {
	value, found := (*hash)[s]
	return value, found
}

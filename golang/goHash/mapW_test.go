package goHash

import (
	"os"
	"path"
	"testing"
)

func TestloadWords(t *testing.T) {
	fp, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	testFile := path.Join(fp, "data", "test.txt")
	mapW := NewMapWord()
	if err := mapW.LoadMapWords(testFile); err != nil {
		t.Error(err)
	}
	if mapW.Length != 3 {
		t.Fail()
	}
	if !mapW.MW.Has("第一个") {
		t.Fail()
	}
	if !mapW.MW.Has("第二个") {
		t.Fail()
	}
	if !mapW.MW.Has("第三个") {
		t.Fail()
	}
	if mapW.MW.Has("第三123个") {
		t.Fail()
	}
}

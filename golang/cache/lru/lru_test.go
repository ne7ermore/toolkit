package utils

import (
	"testing"
)

func TestCache(t *testing.T) {
	cache := GetCache()
	cache.SetCap(2)

	cache.Add("1", 1)
	cache.Add("2", 2)
	res := cache.Get("1")
	if resint, ok := res.(int); !ok {
		t.Errorf("not int")
	} else if resint != 1 {
		t.Errorf("not = 1")
	}

	cache.Add("3", 3)
	res = cache.Get("2")
	if res != nil {
		t.Errorf("should be nil")
	}
	cache.Add("4", 4)
	res = cache.Get("1")
	if res != nil {
		t.Errorf("should be nil")
	}

	res = cache.Get("3")
	if resint, ok := res.(int); !ok {
		t.Errorf("not int")
	} else if resint != 3 {
		t.Errorf("not = 3")
	}

	res = cache.Get("4")
	if resint, ok := res.(int); !ok {
		t.Errorf("not int")
	} else if resint != 4 {
		t.Errorf("not = 4")
	}
}

package common

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	fmt.Println(Hash("你"))
	fmt.Println(Hash("好"))
	fmt.Println(Hash("你好"))
	fmt.Println(Hash("你好1"))
	fmt.Println(Hash("你好2"))
}

func TestShuffle(t *testing.T) {
	a := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println(Shuffle(a, rng))
}

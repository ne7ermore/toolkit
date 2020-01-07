package common

import (
	"hash/fnv"
	"math/rand"
)

func Hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func Shuffle(slice []uint64, r *rand.Rand) (sliceAfter []uint64) {
	for _, v := range r.Perm(len(slice)) {
		sliceAfter = append(sliceAfter, slice[v])
	}
	return sliceAfter
}

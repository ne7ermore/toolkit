package authenticate

import (
	"encoding/hex"
	"math/rand"
	"time"
)

func GenSalt(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(r.Int31()%255 + 1)
	}
	return hex.EncodeToString(b)
}

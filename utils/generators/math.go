package generators

import (
	"math/rand"
	"time"
)

func RandomInt32() int32 {
	rand.Seed(time.Now().UnixNano())

	return int32(rand.Intn(int(2147683647)))
}

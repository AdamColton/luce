package lrand

import (
	"math/rand"
	"unsafe"
)

func Int63() int64 {
	var u int64 = 0
	x := (*[8]byte)(unsafe.Pointer(&u))
	rand.Read(x[:])
	x[7] &= 127
	return u
}

// Source returns a *rand.Rand seeded from crypto/rand
func New() *rand.Rand {
	return rand.New(CryptoSource{})
}

type CryptoSource struct{}

func (CryptoSource) Int63() int64 {
	return Int63()
}

func (CryptoSource) Seed(seed int64) {
	panic("Do not seed CryptoSource")
}

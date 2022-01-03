package lrand

import (
	crand "crypto/rand"
	"math/rand"
	"unsafe"

	"github.com/adamcolton/luce/lerr"
)

// Int63 is used to generate an int64 from crypto/rand.
func Int63() int64 {
	var u int64 = 0
	x := (*[8]byte)(unsafe.Pointer(&u))
	crand.Read(x[:])
	x[7] &= 127
	return u
}

// Source returns a *rand.Rand seeded from crypto/rand
func New() *rand.Rand {
	return rand.New(CryptoSource{})
}

// CryptoSource fulfills rand.Source
type CryptoSource struct{}

// Int63 fulfills rand.Source
func (CryptoSource) Int63() int64 {
	return Int63()
}

// ErrDoNotSeed is the panic value is Seed is called on a CryptoSource.
const ErrDoNotSeed = lerr.Str("Do not seed CryptoSource")

// Seed is required for rand.Source, but should not be used.
func (CryptoSource) Seed(seed int64) {
	panic(ErrDoNotSeed)
}

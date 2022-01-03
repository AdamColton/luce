// Package lrand provides random numbers from crypto/rand, for when they must
// not be predictable.
package lrand

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"

	"github.com/adamcolton/luce/lerr"
)

// read is a random uint64 from crypto/rand. It panics if crypto/rand fails.
func read() uint64 {
	var b [8]byte
	_, err := crand.Read(b[:])
	lerr.Panic(err)
	return binary.BigEndian.Uint64(b[:])
}

// Int63 is used to generate a non-negative int64 from crypto/rand. It panics if
// crypto/rand fails.
func Int63() int64 {
	return int64(read() >> 1)
}

// New returns a *rand.Rand that draws from crypto/rand, through CryptoSource.
func New() *rand.Rand {
	return rand.New(CryptoSource{})
}

// CryptoSource fulfills rand.Source and rand.Source64. It has no state, so it is
// safe to use from several Go routines.
type CryptoSource struct{}

// Int63 fulfills rand.Source.
func (CryptoSource) Int63() int64 {
	return Int63()
}

// Uint64 fulfills rand.Source64, so that a *rand.Rand gets all 64 bits from one
// read rather than building them from two Int63 values.
func (CryptoSource) Uint64() uint64 {
	return read()
}

// ErrDoNotSeed is the panic value if Seed is called on a CryptoSource.
const ErrDoNotSeed = lerr.Str("Do not seed CryptoSource")

// Seed is required for rand.Source, but should not be used. It panics with
// ErrDoNotSeed.
func (CryptoSource) Seed(seed int64) {
	panic(ErrDoNotSeed)
}

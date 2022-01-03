package lrand_test

import (
	"math/rand"
	"testing"

	"github.com/adamcolton/luce/util/lrand"
	"github.com/stretchr/testify/assert"
)

// New is a *rand.Rand, so everything it has works on crypto/rand.
func TestCryptoSource(t *testing.T) {
	r := lrand.New()
	for i := 0; i < 100; i++ {
		n := r.Intn(10)
		assert.True(t, n >= 0 && n < 10)
	}
	assert.Len(t, r.Perm(5), 5)
}

// Int63 is never negative. Over many draws the high bit of a uint64 is seen set
// and not set, so all 64 bits are used.
func TestInt63(t *testing.T) {
	var src rand.Source64 = lrand.CryptoSource{}
	var high, low int
	for i := 0; i < 1000; i++ {
		assert.True(t, lrand.Int63() >= 0)
		assert.True(t, src.Int63() >= 0)
		if src.Uint64()>>63 == 1 {
			high++
		} else {
			low++
		}
	}
	assert.NotZero(t, high)
	assert.NotZero(t, low)
}

func TestSeed(t *testing.T) {
	assert.PanicsWithValue(t, lrand.ErrDoNotSeed, func() {
		lrand.New().Seed(0)
	})
}

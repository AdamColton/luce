package lrand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoSource(t *testing.T) {
	r := New()
	i := r.Intn(10)
	assert.True(t, i >= 0 && i < 10)

	defer func() {
		assert.Equal(t, ErrDoNotSeed, recover())
	}()
	r.Seed(0)
}

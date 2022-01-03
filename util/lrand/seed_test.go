package lrand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoSource(t *testing.T) {
	for i := 0; i < 1000; i++ {
		assert.True(t, Int63() > 0)
	}
}

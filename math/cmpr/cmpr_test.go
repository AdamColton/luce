package cmpr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualWithin(t *testing.T) {
	assert.True(t, Tolerance(0.01).Equal(1, 1.001))
	assert.True(t, Tolerance(0.01).Equal(1.001, 1))
	assert.False(t, Tolerance(0.0001).Equal(1, 1.001))
	assert.False(t, Tolerance(0.0001).Equal(1.001, 1))
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal(1, 1.0+1e-6))
	assert.True(t, Equal(1.0+1e-6, 1))
	assert.False(t, Equal(1, 1.0+1e-4))
	assert.False(t, Equal(1.0+1e-4, 1))
}

func TestZero(t *testing.T) {
	assert.True(t, Zero(0))
	assert.True(t, Zero(1e-6))
	assert.True(t, Zero(-1e-6))
	assert.False(t, Zero(1e-4))
	assert.False(t, Zero(-1e-4))
}

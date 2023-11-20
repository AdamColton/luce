package funcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffStep(t *testing.T) {
	// The step grows with |x| above 1 so it stays above rounding error.
	small, large := diffStep(0.5), diffStep(1e6)
	assert.InDelta(t, 6e-6, small, 1e-6)
	assert.InDelta(t, 6.0, large, 1.0)
	// x+h is exactly h away from x.
	x := 0.1
	h := diffStep(x)
	assert.Equal(t, h, (x+h)-x)
}

package rye_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

// The Ordered decoders work in place, so the input is no longer the ordered
// encoding afterwards.
func TestOrderedDecodeModifiesInput(t *testing.T) {
	b := make([]byte, 4)
	rye.Serialize.Int32Ordered(b, -2)
	assert.Equal(t, []byte{127, 255, 255, 254}, b)

	assert.Equal(t, int32(-2), rye.Deserialize.Int32Ordered(b))
	assert.Equal(t, []byte{254, 255, 255, 255}, b)

	f := make([]byte, 8)
	rye.Serialize.Float64Ordered(f, 1.5)
	assert.Equal(t, float64(1.5), rye.Deserialize.Float64Ordered(f))
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 248, 63}, f)
}

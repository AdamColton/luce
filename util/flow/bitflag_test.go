package flow_test

import (
	"testing"

	"github.com/adamcolton/luce/util/flow"
	"github.com/stretchr/testify/assert"
)

func TestBitFlag(t *testing.T) {
	bf3 := flow.NewFlag(uint16(1 << 3))
	bf5 := flow.NewFlag(uint16(1 << 5))
	var f uint16

	assert.False(t, bf3.Check(f))
	bf3.Set(&f)
	assert.True(t, bf3.Check(f))
	assert.False(t, bf5.Check(f))

	bf5.Set(&f)
	assert.True(t, bf3.Check(f))
	assert.True(t, bf5.Check(f))

	bf3.Clear(&f)
	assert.False(t, bf3.Check(f))
	assert.True(t, bf5.Check(f))

	bf3.Clear(&f)
	assert.False(t, bf3.Check(f))
}

func TestBitFlagMultipleBits(t *testing.T) {
	both := flow.NewFlag(uint8(0b11))
	f := uint8(0b01)
	assert.False(t, both.Check(f), "every bit of the flag must be set")

	f = 0b111
	assert.True(t, both.Check(f))

	both.Clear(&f)
	assert.Equal(t, uint8(0b100), f)
}

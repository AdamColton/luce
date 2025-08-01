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

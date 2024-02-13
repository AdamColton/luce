package numiter_test

import (
	"testing"

	"github.com/adamcolton/luce/math/numiter"
	"github.com/stretchr/testify/assert"
)

func TestInv(t *testing.T) {
	r := numiter.NewRange(0, 10, 2)
	assert.Equal(t, 3, r.Inv(6))
	// Between steps it truncates toward zero.
	assert.Equal(t, 3, r.Inv(7))
	assert.Equal(t, -2, r.Inv(-4))
	// It does not check the range.
	assert.Equal(t, 50, r.Inv(100))

	f := numiter.NewRange(0.0, 1, 0.25)
	assert.Equal(t, 2, f.Inv(0.5))
	assert.Equal(t, 2, f.Inv(0.6))
	assert.Equal(t, -2, f.Inv(-0.6))
	for i := 0; i < f.Len(); i++ {
		assert.Equal(t, i, f.Inv(f.AtIdx(i)))
	}
}

func TestIncludeReachesEnd(t *testing.T) {
	r := numiter.Include(0.0, 1.0, 0.25)
	assert.Equal(t, 5, r.Len())
	assert.Equal(t, 1.0, r.AtIdx(r.Len()-1))

	// When the step does not divide evenly, the last value is past end.
	i := numiter.Include(0, 10, 3)
	assert.Equal(t, 5, i.Len())
	assert.Equal(t, 12, i.AtIdx(i.Len()-1))
}

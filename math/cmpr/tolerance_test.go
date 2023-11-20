package cmpr_test

import (
	"testing"

	"github.com/adamcolton/luce/math/cmpr"
	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	d := float64(cmpr.DefaultTolerance / 10)
	a := 3.1415
	b := a + d
	assert.True(t, cmpr.Equal(a, b))
	assert.True(t, cmpr.Zero(d))
}

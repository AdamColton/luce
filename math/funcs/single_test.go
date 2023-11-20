package funcs_test

import (
	"testing"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/stretchr/testify/assert"
)

func TestSingle(t *testing.T) {
	var s funcs.S = func(x float64) float64 {
		return 0.3*(x*x) - 2.4*x + 7.8
	}

	d := func(x float64) float64 {
		return 0.15*x - 2.4
	}

	steps.Wrap().Iter().For(func(x float64) {
		dx := d(x)
		sdx := s.DPrecise(x, 1e-12)
		assert.InDelta(t, dx, sdx, 1e-6)
	})
}

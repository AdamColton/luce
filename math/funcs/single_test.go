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
		return 0.6*x - 2.4
	}

	steps.Iter().For(func(x float64) {
		dx := d(x)
		sdx := s.D(x)
		assert.InDelta(t, dx, sdx, 1e-8, x)
	})
}

func TestNewton(t *testing.T) {
	var s funcs.S = func(x float64) float64 {
		return 0.3*(x*x) - 2.4*x + 4.2
	}
	x, y := s.NewtonStepper(1).Run(50, 1e-6)
	assert.InDelta(t, 0.0, s(x), 1e-8, x)
	assert.InDelta(t, 0.0, y, 1e-8, y)
}

func TestSecant(t *testing.T) {
	var s funcs.S = func(x float64) float64 {
		return 0.3*(x*x) - 2.4*x + 4.2
	}
	x, y := s.SecantStepper(1).Run(50, 1e-6)
	assert.InDelta(t, 0.0, s(x), 1e-8, x)
	assert.InDelta(t, 0.0, y, 1e-8, y)
}

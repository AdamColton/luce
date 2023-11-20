package funcs_test

import (
	"math"
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

	// One step by hand: s(1) = 2.1 and s'(1) = -1.8.
	dx, y := s.NewtonStep(1)
	assert.InDelta(t, 2.1/1.8, dx, 1e-9)
	assert.InDelta(t, 2.1, y, 1e-12)
}

func TestSecant(t *testing.T) {
	var s funcs.S = func(x float64) float64 {
		return 0.3*(x*x) - 2.4*x + 4.2
	}
	x, y := s.SecantStepper(1).Run(50, 1e-6)
	assert.InDelta(t, 0.0, s(x), 1e-8, x)
	assert.InDelta(t, 0.0, y, 1e-8, y)
}

func TestRunStopsOnNaN(t *testing.T) {
	// A search that hits NaN stops there, and Run returns the best point it
	// saw before.
	points := [][2]float64{{4, 3}, {2, -1}, {1, 2}, {math.NaN(), math.NaN()}, {0, 0}}
	i := 0
	var s funcs.Step = func() (x, y float64) {
		p := points[i]
		i++
		return p[0], p[1]
	}
	x, y := s.Run(10, 1e-6)
	assert.Equal(t, 2.0, x)
	assert.Equal(t, 1.0, y)
	assert.Equal(t, 4, i)
}

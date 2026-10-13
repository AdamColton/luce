package funcs_test

import (
	"math"
	"testing"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/stretchr/testify/assert"
)

func TestDMAgrees(t *testing.T) {
	// A hand-written gradient for f(x, y) = sin(2x)·y, with a common mistake:
	// the chain rule's factor of 2 is missing from ∂f/∂x.
	f := func(x []float64) float64 {
		return math.Sin(2*x[0]) * x[1]
	}
	wrong := funcs.Multi{Ln: 2, M: f, DM: func(x, buf []float64) []float64 {
		return append(buf[:0], math.Cos(2*x[0])*x[1], math.Sin(2*x[0]))
	}}
	right := funcs.Multi{Ln: 2, M: f, DM: func(x, buf []float64) []float64 {
		return append(buf[:0], 2*math.Cos(2*x[0])*x[1], math.Sin(2*x[0]))
	}}

	x := []float64{0.4, 1.5}
	ok, idx := wrong.DMAgrees(x, funcs.DefaultCheckTol)
	assert.False(t, ok)
	assert.Equal(t, 0, idx)

	ok, idx = right.DMAgrees(x, funcs.DefaultCheckTol)
	assert.True(t, ok)
	assert.Equal(t, -1, idx)

	// CheckDM gives the detail: the wrong partial is off by a factor of 2.
	cs := wrong.CheckDM(x, funcs.DefaultCheckTol, nil)
	assert.False(t, cs[0].Pass)
	assert.InDelta(t, 2*cs[0].DM, cs[0].Precise, 1e-9)
	assert.True(t, cs[1].Pass)
	assert.Equal(t, 0.0, cs[1].Precise)
}

func TestCheckDMFallback(t *testing.T) {
	// An energy measured from a large baseline. The baseline swamps the
	// change a central difference's tiny step makes, so rounding error makes
	// NumericDM disagree with the exact gradient. CheckDM falls back to
	// PrecisePartial, which settles it.
	energy := func(baseline float64) funcs.Multi {
		return funcs.Multi{
			Ln: 2,
			M: func(x []float64) float64 {
				return baseline + math.Sin(x[0]) + x[1]*x[1]
			},
			DM: func(x, buf []float64) []float64 {
				return append(buf[:0], math.Cos(x[0]), 2*x[1])
			},
		}
	}
	x := []float64{0.7, 1.3}

	m := energy(1e8)
	cs := m.CheckDM(x, funcs.DefaultCheckTol, make([]funcs.DMCheck, 0, 2))
	for _, c := range cs {
		assert.True(t, c.Pass)
		assert.Greater(t, math.Abs(c.Numeric-c.DM), 1e-5)
		assert.InDelta(t, c.DM, c.Precise, 1e-6)
	}

	// With a larger baseline, even PrecisePartial's error estimate is above
	// the tolerance, so the check can't pass.
	m = energy(1e10)
	for _, c := range m.CheckDM(x, funcs.DefaultCheckTol, nil) {
		assert.False(t, c.Pass)
		assert.Greater(t, c.ErrEst, funcs.DefaultCheckTol)
	}
}

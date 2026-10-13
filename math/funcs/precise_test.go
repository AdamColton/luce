package funcs_test

import (
	"math"
	"testing"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/stretchr/testify/assert"
)

func TestDPrecise(t *testing.T) {
	// e^30x has a large third derivative, which limits the accuracy of a
	// central difference. Ridders' method extrapolates past it.
	var f funcs.S = func(x float64) float64 { return math.Exp(30 * x) }
	want := 30 * math.Exp(15)

	d, errEst := f.DPrecise(0.5)
	assert.InDelta(t, 0, (d-want)/want, 1e-13)
	assert.Less(t, errEst/want, 1e-10)
	assert.Greater(t, math.Abs(f.D(0.5)-want)/want, 1e-9)
}

func TestRiddersH0(t *testing.T) {
	// sin(1e5·x) changes far faster than the default first step, 0.1, and
	// faster than the step S.D uses. With a first step below its scale,
	// Ridders gets it right.
	var f funcs.S = func(x float64) float64 { return math.Sin(1e5 * x) }
	want := 1e5 * math.Cos(1e5*0.3)

	d, errEst := funcs.Ridders{H0: 1e-6}.D(f, 0.3)
	assert.InDelta(t, want, d, 1e-4)
	assert.Less(t, errEst, 1e-3)

	// With the default H0 the answer is wrong, and the error estimate doesn't
	// show it.
	d, errEst = f.DPrecise(0.3)
	assert.Greater(t, math.Abs(d-want), 1e4)
	assert.Less(t, errEst, 100.0)
}

func TestPrecisePartial(t *testing.T) {
	// f(x, y) = x²y, so ∂f/∂x = 2xy and ∂f/∂y = x².
	var f funcs.M = func(x []float64) float64 {
		return x[0] * x[0] * x[1]
	}
	x := []float64{3, 5}
	d, errEst := f.PrecisePartial(x, 0)
	assert.InDelta(t, 30.0, d, 1e-12)
	assert.Less(t, errEst, 1e-9)
	// x is restored.
	assert.Equal(t, []float64{3, 5}, x)

	m := funcs.Multi{Ln: 2, M: f}
	assert.InDeltaSlice(t, []float64{30, 9}, m.PreciseDM(x, nil), 1e-12)

	// A buffer with the capacity is used.
	buf := make([]float64, 0, 2)
	got := funcs.Ridders{H0: 0.01}.DM(f, 2)(x, buf)
	assert.InDeltaSlice(t, []float64{30, 9}, got, 1e-12)
	assert.Equal(t, &buf[:1][0], &got[0])
}

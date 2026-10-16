package funcs_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/stretchr/testify/assert"
)

// parabolaLine is where the parabola (s, s²) meets the line (t, 2-t): the two
// curves are at the same point when x = (s, t) zeroes both equations. They
// cross at (1, 1) and (-2, 4).
func parabolaLine(withJ bool) *funcs.Equations {
	eq := &funcs.Equations{
		In:  2,
		Out: 2,
		R: func(x, buf []float64) []float64 {
			s, t := x[0], x[1]
			return append(buf[:0], s-t, s*s-(2-t))
		},
	}
	if withJ {
		eq.J = func(x, buf []float64) []float64 {
			s := x[0]
			return append(buf[:0],
				1, -1,
				2*s, 1,
			)
		}
	}
	return eq
}

func ExampleEquations_NewtonStepper() {
	eq := parabolaLine(true)
	for _, start := range [][]float64{{0.8, 0.5}, {-1.5, -1}} {
		x, r := eq.NewtonStepper(start).Run(20, 1e-12, nil)
		fmt.Printf("s=%.6f t=%.6f |r|<1e-12: %t\n", x[0], x[1], math.Abs(r[0])+math.Abs(r[1]) < 1e-12)
	}
	// Output:
	// s=1.000000 t=1.000000 |r|<1e-12: true
	// s=-2.000000 t=-2.000000 |r|<1e-12: true
}

func TestNewtonStep(t *testing.T) {
	// One step by hand from (0, 0): r = (0, -2) and J = [[1, -1], [0, 1]], so
	// J·dx = (0, 2) gives dx = (2, 2).
	for _, withJ := range []bool{true, false} {
		dx, r := parabolaLine(withJ).NewtonStep([]float64{0, 0}, nil)
		assert.InDeltaSlice(t, []float64{2, 2}, dx, 1e-9)
		assert.Equal(t, []float64{0, -2}, r)
	}
}

func TestNewtonNumericJ(t *testing.T) {
	// Without J, the stepper estimates it and still finds both crossings.
	eq := parabolaLine(false)
	x, _ := eq.NewtonStepper([]float64{0.8, 0.5}).Run(20, 1e-12, nil)
	assert.InDeltaSlice(t, []float64{1, 1}, x, 1e-9)
	x, _ = eq.NewtonStepper([]float64{-1.5, -1}).Run(20, 1e-12, nil)
	assert.InDeltaSlice(t, []float64{-2, -2}, x, 1e-9)

	want := parabolaLine(true).J([]float64{0.3, 0.7}, nil)
	got := eq.NumericJ([]float64{0.3, 0.7}, nil)
	assert.InDeltaSlice(t, want, got, 1e-9)
}

func TestNewtonSingular(t *testing.T) {
	// The line (t, 0) touches the parabola at the origin without crossing it.
	// At s = 0 the curves' tangents are parallel, so the Jacobian is singular
	// and the step is not finite.
	eq := &funcs.Equations{
		In:  2,
		Out: 2,
		R: func(x, buf []float64) []float64 {
			s, t := x[0], x[1]
			return append(buf[:0], s-t, s*s)
		},
	}
	dx, _ := eq.NewtonStep([]float64{0, 0.5}, nil)
	assert.True(t, math.IsNaN(dx[0]) || math.IsInf(dx[0], 0), dx)
}

func TestNewtonThreeVariables(t *testing.T) {
	// A sphere of radius 3, the plane x = y and the plane z = 1 meet at
	// (2, 2, 1) and (-2, -2, 1).
	eq := &funcs.Equations{
		In:  3,
		Out: 3,
		R: func(v, buf []float64) []float64 {
			x, y, z := v[0], v[1], v[2]
			return append(buf[:0], x*x+y*y+z*z-9, x-y, z-1)
		},
	}
	x, r := eq.NewtonStepper([]float64{1, 3, 0}).Run(50, 1e-12, nil)
	assert.InDeltaSlice(t, []float64{2, 2, 1}, x, 1e-9)
	assert.InDeltaSlice(t, []float64{0, 0, 0}, r, 1e-9)
}

func TestVecStepRun(t *testing.T) {
	// Run keeps the step with the smallest residual and stops at the first
	// step that isn't finite.
	steps := [][2]float64{{4, 3}, {2, -1}, {1, 2}, {math.Inf(1), math.NaN()}, {0, 0}}
	i := 0
	var s funcs.VecStep = func() (x, r []float64) {
		p := steps[i]
		i++
		return []float64{p[0]}, []float64{p[1]}
	}
	buf := make([]float64, 4)
	x, r := s.Run(10, 1e-6, buf)
	assert.Equal(t, []float64{2}, x)
	assert.Equal(t, []float64{-1}, r)
	assert.Equal(t, 4, i)
	assert.Equal(t, 2.0, buf[0], "results are written to buf")
}

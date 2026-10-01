package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/slice"
)

// M is a function of several variables.
type M func([]float64) float64

// PartialDerivative estimates the derivative of the function with respect to
// x[idx] with a central difference, (f(x+h) - f(x-h)) / 2h, where h is
// DiffStep(x[idx]). It calls the function twice. x[idx] is changed while it
// runs and restored before it returns.
func (fn M) PartialDerivative(x []float64, idx int) float64 {
	xi := x[idx]
	h := DiffStep(xi)
	x[idx] = xi + h
	d := fn(x)
	x[idx] = xi - h
	d -= fn(x)
	x[idx] = xi
	return d / (2 * h)
}

// DiffStep is the step for a central difference at x: the cube root of the
// float64 machine epsilon, about 6e-6, scaled by |x| when |x| is above 1. That
// balances the error of the difference formula against rounding error, for
// a relative error near 1e-10 on smooth functions. The step is adjusted so
// that x+h is exactly h away from x.
func DiffStep(x float64) float64 {
	h := cbrtEpsilon * math.Max(1, math.Abs(x))
	return (x + h) - x
}

var cbrtEpsilon = math.Cbrt(0x1p-52)

// DM is the gradient of an M: it returns the partial derivatives at x, using buf
// if it has the capacity.
type DM func(x, buf []float64) []float64

// Multi is a function of Ln variables and, optionally, its gradient.
type Multi struct {
	Ln int
	M  M
	DM DM
}

// NumericDM is a DM that estimates each partial derivative with
// M.PartialDerivative.
func (s *Multi) NumericDM(x, buf []float64) []float64 {
	out := slice.NewBuffer(buf).Slice(s.Ln)
	for i := range out {
		out[i] = s.M.PartialDerivative(x, i)
	}

	return out
}

// GetDM returns DM, or NumericDM if DM is nil.
func (s *Multi) GetDM() DM {
	if s.DM == nil {
		return s.NumericDM
	}
	return s.DM
}

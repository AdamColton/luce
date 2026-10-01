package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/slice"
)

// M is a function of several variables.
type M func([]float64) float64

// == projects.Code.luce.funcs ==
// [ ] PreciseDM
//	A DM using Richardson extrapolation (Ridders' method) for derivatives
//	accurate to about 1e-13, at 10-20 calls each instead of 2. The package is
//	meant as a general tool for multivariate functions, not only descent, so
//	it should also document how reliable PartialDerivative is (its error
//	grows with f''' and with features near the scale of DiffStep) and offer
//	ways to trade precision against cost.
//	https://en.wikipedia.org/wiki/Richardson_extrapolation

// PartialDerivative estimates the derivative of fn with respect to x[idx]
// numerically, to a relative error of about 1e-10 for smooth functions. It
// calls fn twice. It writes to x[idx] during those calls and restores it
// before returning, so x must not be used concurrently.
func (fn M) PartialDerivative(x []float64, idx int) float64 {
	// A central difference: (f(x+h) - f(x-h)) / 2h.
	// https://en.wikipedia.org/wiki/Numerical_differentiation
	xi := x[idx]
	h := DiffStep(xi)
	x[idx] = xi + h
	d := fn(x)
	x[idx] = xi - h
	d -= fn(x)
	x[idx] = xi
	return d / (2 * h)
}

// DiffStep returns the step h to use for a central-difference derivative at x.
func DiffStep(x float64) float64 {
	// The error of a central difference has two parts. The formula's own error
	// shrinks with h², so h should be small. Rounding error in f(x+h) - f(x-h)
	// grows with ε/h, so h shouldn't be too small. Their sum is smallest when
	// h is about ∛ε, which gives a relative error near ε^(2/3), about 1e-10.
	// Scaling by |x| above 1 keeps x+h from rounding back to x when x is large.
	// https://en.wikipedia.org/wiki/Numerical_differentiation#Step_size
	h := cbrtEpsilon * math.Max(1, math.Abs(x))

	// this is not an error
	// x + h may not land exactly h away from x.
	// Recomputing h this way gives the distance actually stepped
	return (x + h) - x
}

// cbrtEpsilon is ∛ε, where ε = 2⁻⁵² is the float64 machine epsilon.
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

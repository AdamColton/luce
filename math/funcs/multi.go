package funcs

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
)

// M is a function of several variables.
type M func([]float64) float64

// PartialDerivative estimates the derivative of the function with respect to
// x[idx]. It halves the step until the difference between the two sides is
// close to zero. x[idx] is changed while it runs and restored before it returns.
func (fn M) PartialDerivative(x []float64, idx int) float64 {
	// d is the delta between x+step and x-step
	// first we're looking for a step size that gets d close to zero
	const small cmpr.Tolerance = 1e-3
	step := 1e-2
	d := 1.0
	xi := x[idx]
	for !small.Zero(d) && step > 1e-12 {
		step /= 2
		x[idx] = xi + step
		d = fn(x)
		x[idx] = xi - step
		d -= fn(x)
		x[idx] = xi
	}
	return d / (2 * step)
}

// DM is the gradient of an M: it returns the partial derivatives at x, using buf
// if it has the capacity.
type DM func(x, buf []float64) []float64

// Multi is a function of Ln variables and, optionally, its gradient.
type Multi struct {
	Ln int
	M  M
	DM DM
}

// AnalyticDM is a DM that estimates each partial derivative with
// M.PartialDerivative. Despite the name it is numeric.
func (s *Multi) AnalyticDM(x, buf []float64) []float64 {
	out := slice.NewBuffer(buf).Slice(s.Ln)
	for i := range out {
		out[i] = s.M.PartialDerivative(x, i)
	}

	return out
}

// GetDM returns DM, or AnalyticDM if DM is nil.
func (s *Multi) GetDM() DM {
	if s.DM == nil {
		return s.AnalyticDM
	}
	return s.DM
}

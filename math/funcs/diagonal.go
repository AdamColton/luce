package funcs

import (
	"math"
)

// Diagonal steps down the gradient with each variable scaled by how sharply f
// curves along it. That handles variables on very different scales (volts and
// amps), but not variables that depend on each other. It costs 2 extra calls
// to M per variable at each step.
type Diagonal struct {
	// h is the diagonal of the Hessian at the current point.
	h []float64
}

// Reset implements Direction.
func (dg *Diagonal) Reset(n int) {
	if cap(dg.h) < n {
		dg.h = make([]float64, n)
	}
	dg.h = dg.h[:n]
}

// fourthRootEpsilon is the step for a second difference. It balances the
// formula's error against rounding error, as diffStep does for a first
// difference: rounding error in a second difference grows with ε/h², which
// moves the best step up to about ⁴√ε.
var fourthRootEpsilon = math.Pow(0x1p-52, 0.25)

// minCurvature keeps Diagonal from dividing by a curvature of zero.
const minCurvature = 1e-12

// Dir implements Direction.
func (dg *Diagonal) Dir(m M, x []float64, f float64, g, d []float64) float64 {
	// Newton's method steps by -g/f'' in one variable, which lands exactly on
	// the minimum of a quadratic. Diagonal does that for each variable on its
	// own, using the diagonal of the Hessian (the second derivatives ∂²f/∂xᵢ²)
	// and ignoring the rest of it.
	// https://en.wikipedia.org/wiki/Newton%27s_method_in_optimization
	//
	// Each second derivative is a second difference:
	// (f(x+h) + f(x-h) - 2f(x)) / h².
	// https://en.wikipedia.org/wiki/Finite_difference#Higher-order_differences
	dg.Reset(len(x))
	for i, xi := range x {
		h := fourthRootEpsilon * math.Max(1, math.Abs(xi))
		x[i] = xi + h
		fp := m(x)
		x[i] = xi - h
		fm := m(x)
		x[i] = xi
		dg.h[i] = (fp + fm - 2*f) / (h * h)
		// The absolute value keeps the step pointing down where f curves
		// downward, which a plain Newton step would turn uphill.
		d[i] = -g[i] / math.Max(math.Abs(dg.h[i]), minCurvature)
	}
	return 1
}

// Update implements Direction.
func (dg *Diagonal) Update(s, y []float64, t float64) {}

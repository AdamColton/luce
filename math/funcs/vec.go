package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/slice"
)

// == projects.Code.luce.funcs ==
// [ ] move vector helpers to math/vec
//	neg, dot, axpy and hasNaN belong in a math/vec package, along with the
//	other vector loops in math/funcs.

// neg returns -a, using buf if it has the capacity. buf may be a itself.
func neg(a, buf []float64) []float64 {
	out := slice.NewBuffer(buf).Slice(len(a))
	for i, ai := range a {
		out[i] = -ai
	}
	return out
}

// dot returns the dot product of a and b.
func dot(a, b []float64) (sum float64) {
	for i, ai := range a {
		sum += ai * b[i]
	}
	return
}

// hasNaN reports whether any element of xs is NaN.
func hasNaN(xs []float64) bool {
	for _, x := range xs {
		if math.IsNaN(x) {
			return true
		}
	}
	return false
}

// axpy adds a*x to y.
func axpy(a float64, x, y []float64) {
	for i, xi := range x {
		y[i] += a * xi
	}
}

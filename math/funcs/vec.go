package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
)

// == projects.Code.luce.funcs ==
// [ ] move vector helpers to math/vec
//	neg, dot, axpy, hasNaN, finite, maxAbs, allClose and solve belong in a
//	math/vec package, along with the other vector loops in math/funcs.

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

// finite reports whether every element of xs is finite: neither NaN nor Inf.
func finite(xs []float64) bool {
	for _, x := range xs {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return false
		}
	}
	return true
}

// maxAbs returns the largest absolute value in xs.
func maxAbs(xs []float64) (m float64) {
	for _, x := range xs {
		m = math.Max(m, math.Abs(x))
	}
	return
}

// allClose reports whether each element of a is within small of the same
// element of b.
func allClose(small cmpr.Tolerance, a, b []float64) bool {
	for i, ai := range a {
		if !small.Equal(ai, b[i]) {
			return false
		}
	}
	return true
}

// solve solves a·x = b for x, where a is an n×n matrix stored row-major and b
// has length n. It writes x into b and uses a as scratch space. If a is
// singular, x holds Inf or NaN.
func solve(a, b []float64) {
	// Gaussian elimination with partial pivoting: at each column, swap up the
	// row with the largest entry, so no row is scaled by a large factor.
	// https://en.wikipedia.org/wiki/Gaussian_elimination#Partial_pivoting
	n := len(b)
	for c := 0; c < n; c++ {
		p := c
		for i := c + 1; i < n; i++ {
			if math.Abs(a[i*n+c]) > math.Abs(a[p*n+c]) {
				p = i
			}
		}
		if p != c {
			for j := c; j < n; j++ {
				a[c*n+j], a[p*n+j] = a[p*n+j], a[c*n+j]
			}
			b[c], b[p] = b[p], b[c]
		}
		for i := c + 1; i < n; i++ {
			f := a[i*n+c] / a[c*n+c]
			for j := c; j < n; j++ {
				a[i*n+j] -= f * a[c*n+j]
			}
			b[i] -= f * b[c]
		}
	}
	// Back substitution, from the last row up.
	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			b[i] -= a[i*n+j] * b[j]
		}
		b[i] /= a[i*n+i]
	}
}

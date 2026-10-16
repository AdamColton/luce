package funcs

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
)

// R is a system of equations, each written so that it is zero at a solution.
// It returns the value of every equation at x, using buf if it has the
// capacity.
type R func(x, buf []float64) []float64

// J is the Jacobian of an R: the partial derivative of each equation with
// respect to each variable at x. It is row-major, so row i holds equation i,
// and uses buf if it has the capacity.
type J func(x, buf []float64) []float64

// Equations is a system of Out equations in In unknowns and, optionally, its
// Jacobian.
type Equations struct {
	In, Out int
	R       R
	J       J
}

// NumericJ is a J that estimates each partial derivative by a central
// difference, as M.PartialDerivative does. It calls R 2·In times, and writes to
// x during those calls, restoring it before returning, so x must not be used
// concurrently.
func (eq *Equations) NumericJ(x, buf []float64) []float64 {
	n, m := eq.In, eq.Out
	// The residuals at each probe are written to the m slots after the
	// Jacobian, so the caller's buf covers them as well.
	out := slice.NewBuffer(buf).Slice(m * (n + 1))
	jac, probe := out[:m*n], out[m*n:]
	for j, xj := range x[:n] {
		h := diffStep(xj)
		x[j] = xj + h
		r := eq.R(x, probe)
		for i, ri := range r {
			jac[i*n+j] = ri
		}
		x[j] = xj - h
		r = eq.R(x, probe)
		for i, ri := range r {
			jac[i*n+j] = (jac[i*n+j] - ri) / (2 * h)
		}
		x[j] = xj
	}
	return jac
}

// GetJ returns J, or NumericJ if J is nil.
func (eq *Equations) GetJ() J {
	if eq.J == nil {
		return eq.NumericJ
	}
	return eq.J
}

// == projects.Code.luce.funcs ==
// [ ] Newton for non-square systems
//  NewtonStep needs as many equations as unknowns. With more equations, the
//  step is the least-squares solution (Gauss-Newton), which is where
//  Levenberg-Marquardt starts.

// NewtonStep returns the step dx that Newton's method takes from x, and r, the
// value of the equations at x. It requires In == Out. When the Jacobian at x is
// singular, for example where two curves touch without crossing, dx holds Inf
// or NaN. buf can hold dx, r and the Jacobian together; the results share its
// memory.
func (eq *Equations) NewtonStep(x, buf []float64) (dx, r []float64) {
	n := eq.In
	all := slice.NewBuffer(buf).Slice(n * (n + 3))
	dx, r = all[:n:n], all[n:2*n:2*n]
	r = eq.R(x, r)
	return eq.newtonStepR(x, r, dx, all[2*n:]), r
}

// newtonStepR is NewtonStep when r, the value at x, is already known. jbuf
// must hold n·(n+1) values, for NumericJ.
func (eq *Equations) newtonStepR(x, r, dx, jbuf []float64) []float64 {
	// Near x, R(x+dx) ≈ R(x) + J·dx, so the step that brings that to zero
	// solves J·dx = -R(x).
	// https://en.wikipedia.org/wiki/Newton%27s_method#Systems_of_equations
	jac := eq.GetJ()(x, jbuf)
	neg(r, dx)
	solve(jac, dx)
	return dx
}

// VecStep calls the next step of a search in several variables and returns
// the x it reached and the value r of the equations there. The slices may be
// overwritten by the next call.
type VecStep func() (x, r []float64)

// NewtonStepper returns a VecStep that applies Newton's method, starting at x,
// to look for a solution. It updates x in place. It requires In == Out.
func (eq *Equations) NewtonStepper(x []float64) VecStep {
	n := eq.In
	all := make([]float64, n*(n+3))
	r, dx, jbuf := all[:n:n], all[n:2*n:2*n], all[2*n:]
	r = eq.R(x, r)
	return func() ([]float64, []float64) {
		eq.newtonStepR(x, r, dx, jbuf)
		axpy(1, dx, x)
		r = eq.R(x, r)
		return x, r
	}
}

// Run takes up to max steps and returns the x and r with the smallest r, by
// its largest absolute value, that it saw. It stops early when every component
// of both x and r changes by less than small, or when a step returns NaN or
// Inf. buf can hold both results and a copy of each; the results share its
// memory.
func (s VecStep) Run(max int, small cmpr.Tolerance, buf []float64) (x, r []float64) {
	x, r = s()
	n, m := len(x), len(r)
	all := slice.NewBuffer(buf).Slice(2 * (n + m))
	bx, br := all[:n:n], all[n:n+m:n+m]
	px, pr := all[n+m:2*n+m:2*n+m], all[2*n+m:]
	copy(bx, x)
	copy(br, r)
	best := maxAbs(r)
	for i := 0; i < max; i++ {
		copy(px, x)
		copy(pr, r)
		x, r = s()
		// A singular Jacobian gives a step of Inf or NaN, which can't lead
		// anywhere useful.
		if !finite(x) || !finite(r) {
			break
		}
		if m := maxAbs(r); m < best {
			best = m
			copy(bx, x)
			copy(br, r)
		}
		if allClose(small, px, x) && allClose(small, pr, r) {
			break
		}
	}
	return bx, br
}

package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/slice"
)

// DefaultH0 is the H0 that a zero Ridders uses.
const DefaultH0 = 0.1

// Ridders estimates derivatives with Ridders' method, to a relative error near
// 1e-13 for smooth functions, along with an estimate of that error. Each
// derivative costs up to 20 calls, against 2 for S.D and M.PartialDerivative.
// The zero value uses the defaults.
//
// The central differences that S.D and M.PartialDerivative use are accurate to
// about 1e-10 when f is smooth on the scale of their step, which is about
// 6e-6·max(1, |x|). Where f changes much faster than that, for example
// sin(1e5·x), or has a large third derivative, they can be badly off. Ridders
// starts with a larger step and shrinks it, so it handles the second case. For
// the first, set H0 below the scale on which f changes. With an H0 that is too
// large, the estimate can be wrong while the error estimate is small: for
// sin(1e5·x), the default H0 gives an answer off by 100% with an error
// estimate of 0.02%. Where f isn't differentiable, no numeric method gives a
// meaningful answer.
type Ridders struct {
	// H0 is the first step, relative to max(1, |x|). The steps shrink from
	// there. 0 means DefaultH0.
	H0 float64
}

const (
	// riddersShrink is the factor each step is smaller than the last.
	riddersShrink = 1.4
	// riddersSteps is the most steps Ridders takes, at 2 calls each.
	riddersSteps = 10
	// riddersSafe stops the search once the error of the newest estimate is
	// this many times the best error so far: rounding has taken over.
	riddersSafe = 2.0
)

// D estimates the derivative of fn at x, and the error of that estimate.
func (r Ridders) D(fn S, x float64) (d, errEst float64) {
	h0 := r.h0() * math.Max(1, math.Abs(x))
	return ridders(fn, x, h0)
}

// Partial estimates the derivative of fn with respect to x[idx], and the error
// of that estimate. It writes to x[idx] during the calls and restores it before
// returning, so x must not be used concurrently.
func (r Ridders) Partial(fn M, x []float64, idx int) (d, errEst float64) {
	xi := x[idx]
	h0 := r.h0() * math.Max(1, math.Abs(xi))
	d, errEst = ridders(func(v float64) float64 {
		x[idx] = v
		return fn(x)
	}, xi, h0)
	x[idx] = xi
	return
}

// DM returns a DM for fn of n variables that estimates each partial derivative
// with Partial.
func (r Ridders) DM(fn M, n int) DM {
	return func(x, buf []float64) []float64 {
		out := slice.NewBuffer(buf).Slice(n)
		for i := range out {
			out[i], _ = r.Partial(fn, x, i)
		}
		return out
	}
}

func (r Ridders) h0() float64 {
	if r.H0 == 0 {
		return DefaultH0
	}
	return r.H0
}

// ridders estimates the derivative of fn at x, starting from the step h.
func ridders(fn func(float64) float64, x, h float64) (d, errEst float64) {
	// Ridders' method: take central differences at steps h, h/1.4, h/1.4²...
	// and extrapolate them to a step of 0. A central difference's error is a
	// series in h², so a pair of estimates at two steps can be combined to
	// cancel the h² term, a pair of those to cancel h⁴, and so on: Neville's
	// algorithm, applied to Richardson extrapolation. a[j][i] is the estimate
	// from step i with j terms cancelled. How much neighbouring entries
	// disagree estimates the error, and the search keeps the entry where that
	// is smallest. It stops when the newest row is much worse than the best,
	// because the step has become small enough for rounding to dominate.
	// https://en.wikipedia.org/wiki/Ridders%27_method
	// Numerical Recipes in C, 2nd ed., §5.7, dfridr.
	var a [riddersSteps][riddersSteps]float64
	central := func(h float64) float64 {
		return (fn(x+h) - fn(x-h)) / (2 * h)
	}
	a[0][0] = central(h)
	d, errEst = a[0][0], math.Inf(1)
	const shrink2 = riddersShrink * riddersShrink
	for i := 1; i < riddersSteps; i++ {
		h /= riddersShrink
		a[0][i] = central(h)
		fac := shrink2
		for j := 1; j <= i; j++ {
			a[j][i] = (a[j-1][i]*fac - a[j-1][i-1]) / (fac - 1)
			fac *= shrink2
			e := math.Max(math.Abs(a[j][i]-a[j-1][i]), math.Abs(a[j][i]-a[j-1][i-1]))
			if e <= errEst {
				d, errEst = a[j][i], e
			}
		}
		if math.Abs(a[i][i]-a[i-1][i-1]) >= riddersSafe*errEst {
			break
		}
	}
	return
}

// DPrecise estimates the derivative of fn at x with Ridders' method, and the
// error of that estimate. It costs up to 20 calls; see Ridders.
func (fn S) DPrecise(x float64) (d, errEst float64) {
	return Ridders{}.D(fn, x)
}

// PrecisePartial estimates the derivative of fn with respect to x[idx] with
// Ridders' method, and the error of that estimate. It costs up to 20 calls;
// see Ridders. It writes to x[idx] during the calls and restores it before
// returning, so x must not be used concurrently.
func (fn M) PrecisePartial(x []float64, idx int) (d, errEst float64) {
	return Ridders{}.Partial(fn, x, idx)
}

// PreciseDM is a DM that estimates each partial derivative with
// M.PrecisePartial.
func (s *Multi) PreciseDM(x, buf []float64) []float64 {
	return Ridders{}.DM(s.M, s.Ln)(x, buf)
}

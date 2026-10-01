package funcs

import (
	"math"
)

// Direction picks the direction of each step of a Descender. A Direction can
// keep state between steps, so each Descender needs its own.
type Direction interface {
	// Reset starts again for n variables, forgetting any history.
	Reset(n int)
	// Dir writes the direction to step from x into d, where m has the value f
	// and the gradient g, and returns the step length for the line search to
	// try first. The direction should lead down (d·g < 0); if it doesn't, the
	// Descender uses -g and calls Reset. Dir may call m, and change x while
	// it does, as long as x is restored before Dir returns. Those calls count
	// as evaluations.
	Dir(m M, x []float64, f float64, g, d []float64) (t float64)
	// Update is told about each step the line search accepted: s is the
	// change in x, y the change in the gradient and t the step length used.
	Update(s, y []float64, t float64)
}

// Gradient steps straight down the gradient: plain gradient descent. It is
// simple, but slow when the variables are on very different scales or depend
// on each other.
//
// https://en.wikipedia.org/wiki/Gradient_descent
type Gradient struct {
	// t is the last step length the line search accepted.
	t float64
}

// Reset implements Direction.
func (gr *Gradient) Reset(n int) {
	gr.t = 0
}

// Dir implements Direction.
func (gr *Gradient) Dir(m M, x []float64, f float64, g, d []float64) float64 {
	for i, gi := range g {
		d[i] = -gi
	}
	if gr.t == 0 {
		return 1
	}
	// The line search can only shrink a step, so starting from the last
	// accepted length would let the step shrink but never grow. Starting at
	// four times it lets the step grow again where the function allows it.
	return 4 * gr.t
}

// Update implements Direction.
func (gr *Gradient) Update(s, y []float64, t float64) {
	gr.t = t
}

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

// DefaultMemory is the Memory LBFGS uses when it is 0.
const DefaultMemory = 6

// LBFGS is limited-memory BFGS, a quasi-Newton method. It learns how f
// curves from the last Memory steps, so it handles both variables on very
// different scales and variables that depend on each other, without extra
// calls to M. It is the Descender's default.
//
// https://en.wikipedia.org/wiki/Limited-memory_BFGS
type LBFGS struct {
	// Memory is how many recent steps it remembers. 0 means DefaultMemory.
	Memory int

	// s[i] and y[i] are the change in x and in the gradient over a remembered
	// step, oldest first, and rho[i] = 1/(y[i]·s[i]).
	s, y [][]float64
	rho  []float64
	// alpha is scratch space for Dir.
	alpha []float64
}

// Reset implements Direction.
func (l *LBFGS) Reset(n int) {
	l.s, l.y, l.rho = l.s[:0], l.y[:0], l.rho[:0]
}

func (l *LBFGS) memory() int {
	if l.Memory == 0 {
		return DefaultMemory
	}
	return l.Memory
}

// Dir implements Direction.
func (l *LBFGS) Dir(m M, x []float64, f float64, g, d []float64) float64 {
	// Newton's method would step by -H⁻¹g, where H is the Hessian. BFGS builds
	// an estimate of H⁻¹ from the steps so far: each step s and the change in
	// the gradient y over it satisfy H s ≈ y. L-BFGS never forms the matrix.
	// It applies the estimate to g directly from the last few (s, y) pairs
	// with the "two-loop recursion".
	// https://en.wikipedia.org/wiki/Limited-memory_BFGS#Algorithm
	copy(d, g)
	k := len(l.s)
	if k == 0 {
		// No history yet, so no curvature estimate: step down the gradient,
		// scaled so that no variable moves by more than 1 on the first try.
		scale := 0.0
		for _, gi := range g {
			scale = math.Max(scale, math.Abs(gi))
		}
		if scale == 0 {
			scale = 1
		}
		for i := range d {
			d[i] = -d[i] / scale
		}
		return 1
	}

	if cap(l.alpha) < k {
		l.alpha = make([]float64, k)
	}
	alpha := l.alpha[:k]
	// First loop, newest pair to oldest.
	for i := k - 1; i >= 0; i-- {
		alpha[i] = l.rho[i] * dot(l.s[i], d)
		axpy(-alpha[i], l.y[i], d)
	}
	// The starting estimate of H⁻¹ is γI, with γ = s·y / y·y from the newest
	// pair: the inverse curvature along the last step. This is what scales the
	// step, and why the line search usually accepts t = 1.
	last := k - 1
	gamma := dot(l.s[last], l.y[last]) / dot(l.y[last], l.y[last])
	for i := range d {
		d[i] *= gamma
	}
	// Second loop, oldest pair to newest.
	for i := 0; i < k; i++ {
		beta := l.rho[i] * dot(l.y[i], d)
		axpy(alpha[i]-beta, l.s[i], d)
	}
	// d now holds H⁻¹g; the step goes the opposite way.
	for i := range d {
		d[i] = -d[i]
	}
	return 1
}

// Update implements Direction.
func (l *LBFGS) Update(s, y []float64, t float64) {
	// Only remember a step where f curved upward along it (s·y > 0, the
	// "curvature condition"). That keeps the H⁻¹ estimate positive definite,
	// which guarantees that Dir points down. The 1e-12 margin, relative to the
	// sizes of s and y, also skips steps where rounding hides the sign.
	sy := dot(s, y)
	if !(sy > 1e-12*math.Sqrt(dot(s, s)*dot(y, y))) {
		return
	}
	var sBuf, yBuf []float64
	if len(l.s) == l.memory() {
		// Full: drop the oldest pair and reuse its slices for the new one.
		sBuf, yBuf = l.s[0], l.y[0]
		copy(l.s, l.s[1:])
		copy(l.y, l.y[1:])
		copy(l.rho, l.rho[1:])
		l.s, l.y, l.rho = l.s[:len(l.s)-1], l.y[:len(l.y)-1], l.rho[:len(l.rho)-1]
	}
	l.s = append(l.s, append(sBuf[:0], s...))
	l.y = append(l.y, append(yBuf[:0], y...))
	l.rho = append(l.rho, 1/sy)
}

// axpy adds a*x to y.
func axpy(a float64, x, y []float64) {
	for i, xi := range x {
		y[i] += a * xi
	}
}

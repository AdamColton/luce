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
	// try first. It must not change x. m can be called to learn more about the
	// function; those calls count as evaluations.
	Dir(m M, x []float64, f float64, g, d []float64) (t float64)
	// Update is told about each step the line search accepted: s is the
	// change in x, y the change in the gradient and t the step length used.
	Update(s, y []float64, t float64)
}

// Gradient steps straight down the gradient. It is the plain form of
// gradient descent: simple, but slow when the variables are on very
// different scales. The first step length tried is four times the last one
// accepted, so the step can grow back after the line search shrinks it.
type Gradient struct {
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
	return 4 * gr.t
}

// Update implements Direction.
func (gr *Gradient) Update(s, y []float64, t float64) {
	gr.t = t
}

// Diagonal steps down the gradient with each component divided by the
// curvature of f along that variable, the diagonal of the Hessian. That
// rescales the variables, so it handles variables on different scales (volts
// and amps), but not variables that are strongly coupled. The curvature is
// estimated with second differences, which costs 2 calls to M per variable at
// each step.
type Diagonal struct {
	h []float64
}

// Reset implements Direction.
func (dg *Diagonal) Reset(n int) {
	if cap(dg.h) < n {
		dg.h = make([]float64, n)
	}
	dg.h = dg.h[:n]
}

// fourthRootEpsilon is the step for a second difference, which balances
// error the same way DiffStep does for a first difference.
var fourthRootEpsilon = math.Pow(0x1p-52, 0.25)

// minCurvature keeps Diagonal from dividing by a curvature of zero.
const minCurvature = 1e-12

// Dir implements Direction.
func (dg *Diagonal) Dir(m M, x []float64, f float64, g, d []float64) float64 {
	dg.Reset(len(x))
	for i, xi := range x {
		h := fourthRootEpsilon * math.Max(1, math.Abs(xi))
		x[i] = xi + h
		fp := m(x)
		x[i] = xi - h
		fm := m(x)
		x[i] = xi
		dg.h[i] = (fp + fm - 2*f) / (h * h)
		d[i] = -g[i] / math.Max(math.Abs(dg.h[i]), minCurvature)
	}
	return 1
}

// Update implements Direction.
func (dg *Diagonal) Update(s, y []float64, t float64) {}

// DefaultMemory is the Memory LBFGS uses when it is 0.
const DefaultMemory = 6

// LBFGS is the limited-memory BFGS method, a quasi-Newton method. It
// estimates the curvature of f from how the gradient changed over the last
// Memory steps, so it costs no extra calls to M and handles both variables on
// different scales and variables that are coupled. It is the Descender's
// default.
type LBFGS struct {
	// Memory is how many recent steps it remembers. 0 means DefaultMemory.
	Memory int
	s, y   [][]float64
	rho    []float64
	alpha  []float64
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
	copy(d, g)
	k := len(l.s)
	if k == 0 {
		// No history yet: a gradient step scaled so the first step changes no
		// variable by more than 1.
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

	// The two-loop recursion.
	if cap(l.alpha) < k {
		l.alpha = make([]float64, k)
	}
	alpha := l.alpha[:k]
	for i := k - 1; i >= 0; i-- {
		alpha[i] = l.rho[i] * dot(l.s[i], d)
		axpy(-alpha[i], l.y[i], d)
	}
	last := k - 1
	gamma := dot(l.s[last], l.y[last]) / dot(l.y[last], l.y[last])
	for i := range d {
		d[i] *= gamma
	}
	for i := 0; i < k; i++ {
		beta := l.rho[i] * dot(l.y[i], d)
		axpy(alpha[i]-beta, l.s[i], d)
	}
	for i := range d {
		d[i] = -d[i]
	}
	return 1
}

// Update implements Direction. A step where the gradient didn't change in a
// way that shows positive curvature is skipped, which keeps the curvature
// estimate positive.
func (l *LBFGS) Update(s, y []float64, t float64) {
	sy := dot(s, y)
	if !(sy > 1e-12*math.Sqrt(dot(s, s)*dot(y, y))) {
		return
	}
	var sBuf, yBuf []float64
	if len(l.s) == l.memory() {
		// Reuse the oldest entry's slices.
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

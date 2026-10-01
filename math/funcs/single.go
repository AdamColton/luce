package funcs

import (
	"math"

	"github.com/adamcolton/luce/math/cmpr"
)

// S is a function of a single variable.
type S func(float64) float64

// DPrecise estimates the derivative at x. It halves the step, starting at
// small, until the difference between the two sides is within small of zero.
func (fn S) DPrecise(x float64, small cmpr.Tolerance) float64 {
	step := float64(small)
	d := 1.0
	for !small.Zero(d) {
		step /= 2
		d = fn(x+step) - fn(x-step)
	}
	return d / (2 * float64(step))
}

// D estimates the derivative of fn at x numerically, to a relative error of
// about 1e-10 for smooth functions. It calls fn twice.
func (fn S) D(x float64) float64 {
	// A central difference, as in M.PartialDerivative.
	h := DiffStep(x)
	return (fn(x+h) - fn(x-h)) / (2 * h)
}

// NewtonStep returns the step dx that Newton's method takes from x, and y, the
// value of the function at x.
func (fn S) NewtonStep(x float64) (dx, y float64) {
	y = fn(x)
	return fn.NewtonStepY(x, y), y
}

// NewtonStepY is NewtonStep when y, the value at x, is already known.
func (fn S) NewtonStepY(x, y float64) (dx float64) {
	m := fn.D(x)
	dx = -y / m

	return
}

// NewtonStepper returns a Step that applies Newton's method, starting at x, to
// look for a zero.
func (fn S) NewtonStepper(x float64) Step {
	y := fn(x)
	return func() (float64, float64) {
		x += fn.NewtonStepY(x, y)
		y = fn(x)
		return x, y
	}
}

// SecantStep returns where the line through (x0, y0) and (x1, y1) crosses zero.
func SecantStep(x0, x1, y0, y1 float64) float64 {
	return x0 - (y0*(x0-x1))/(y0-y1)
}

// SecantStepper returns a Step that applies the secant method, starting at x0, to
// look for a zero. The second point comes from one Newton step.
func (fn S) SecantStepper(x0 float64) Step {
	y0 := fn(x0)
	x1 := x0 + fn.NewtonStepY(x0, y0)
	y1 := fn(x1)
	return func() (float64, float64) {
		x0, x1 = SecantStep(x0, x1, y0, y1), x0
		y0, y1 = fn(x0), y0
		return x0, y0
	}
}

// Step calls the next step of a search and returns the x it reached and the
// value y of the function there.
type Step func() (x, y float64)

type best struct {
	x, y float64
}

func (b *best) update(x, y float64, init bool) {
	y = math.Abs(y)
	if init || y < b.y {
		b.x = x
		b.y = y
	}
}

// Run takes up to max steps and returns the x and y with the smallest absolute
// y that it saw. It stops early when both x and y change by less than small.
func (s Step) Run(max int, small cmpr.Tolerance) (x, y float64) {
	x, y = s()
	b := &best{x, y}
	for i := 0; i < max; i++ {
		cx, cy := s()
		b.update(cx, cy, false)
		dx, dy := x-cx, y-cy
		if small.Zero(dx) && small.Zero(dy) {
			break
		}
		x, y = cx, cy
	}
	return b.x, b.y
}

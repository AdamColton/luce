package funcs

import (
	"math"

	"github.com/adamcolton/luce/math/cmpr"
)

type S func(float64) float64

func (fn S) DPrecise(x float64, small cmpr.Tolerance) float64 {
	step := float64(small)
	d := 1.0
	for !small.Zero(d) {
		step /= 2
		d = fn(x+step) - fn(x-step)
	}
	return d / (2 * float64(step))
}

func (fn S) D(x float64) float64 {
	return fn.DPrecise(x, 1e-6)
}

func (fn S) NewtonStep(x float64) (dx, y float64) {
	y = fn(x)
	return fn.NewtonStepY(x, y), y
}

func (fn S) NewtonStepY(x, y float64) (dx float64) {
	m := fn.D(x)
	dx = -y / m

	return
}

func (fn S) NewtonStepper(x float64) Step {
	y := fn(x)
	return func() (float64, float64) {
		x += fn.NewtonStepY(x, y)
		y = fn(x)
		return x, y
	}
}

func SecantStep(x0, x1, y0, y1 float64) float64 {
	return x0 - (y0*(x0-x1))/(y0-y1)
}

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

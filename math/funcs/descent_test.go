package funcs_test

import (
	"math"
	"testing"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/stretchr/testify/assert"
)

func TestGradientDescent(t *testing.T) {
	// Find where the two circles of the other tests come closest. They cross,
	// so the distance goes to zero.
	d := funcs.NewDescender(funcs.Multi{Ln: 2, M: dist}, []float64{0.4, 0.3})
	r := d.Run()

	assert.True(t, r.Reason.Converged(), r.Reason.String())
	assert.InDelta(t, 0.0, r.F, 1e-6)
	assert.InDelta(t, 0.0, dist(r.X), 1e-6)
	assert.Equal(t, d.Result(), r)
}

// bowl is f(x, y) = (x-1)² + 10(y+2)², with its minimum at (1, -2).
var bowl funcs.M = func(x []float64) float64 {
	a, b := x[0]-1, x[1]+2
	return a*a + 10*b*b
}

func TestDescenderStepByStep(t *testing.T) {
	// Step can drive the descent one step at a time, for example to draw its
	// path. Record sees the start and every step.
	var path []funcs.StepRecord
	d := funcs.NewDescender(funcs.Multi{Ln: 2, M: bowl}, []float64{0, 0})
	d.Record = func(r funcs.StepRecord) {
		path = append(path, r)
	}
	for d.Step() {
		assert.Equal(t, funcs.Running, d.Result().Reason)
	}
	r := d.Result()

	assert.Equal(t, funcs.StepSmall, r.Reason)
	assert.InDelta(t, 1.0, r.X[0], 1e-6)
	assert.InDelta(t, -2.0, r.X[1], 1e-6)
	// Stepping a stopped Descender does nothing.
	assert.False(t, d.Step())
	assert.Equal(t, r, d.Result())

	if assert.Len(t, path, r.Steps+1) {
		start := path[0]
		assert.Equal(t, 0, start.Step)
		assert.Equal(t, []float64{0, 0}, start.X)
		assert.Equal(t, 41.0, start.F)
		assert.Nil(t, start.Dir)
		last := path[len(path)-1]
		assert.Equal(t, r.Steps, last.Step)
		assert.Equal(t, r.X, last.X)
		assert.Len(t, last.Dir, 2)
		assert.Greater(t, last.T, 0.0)
	}
}

func TestDirections(t *testing.T) {
	for name, dir := range map[string]funcs.Direction{
		"gradient": &funcs.Gradient{},
		"diagonal": &funcs.Diagonal{},
		"lbfgs":    &funcs.LBFGS{Memory: 2},
	} {
		d := funcs.NewDescender(funcs.Multi{Ln: 2, M: bowl}, []float64{0, 0})
		d.Direction = dir
		r := d.Run()
		assert.True(t, r.Reason.Converged(), name)
		assert.InDelta(t, 1.0, r.X[0], 1e-6, name)
		assert.InDelta(t, -2.0, r.X[1], 1e-6, name)
	}
}

func TestStop(t *testing.T) {
	newBowl := func(stop funcs.Stop) funcs.Result {
		d := funcs.NewDescender(funcs.Multi{Ln: 2, M: bowl}, []float64{0, 0})
		d.Direction = &funcs.Gradient{}
		d.Stop = stop
		return d.Run()
	}

	r := newBowl(funcs.Stop{MaxSteps: 3})
	assert.Equal(t, funcs.MaxSteps, r.Reason)
	assert.Equal(t, 3, r.Steps)
	assert.False(t, r.Reason.Converged())

	r = newBowl(funcs.Stop{MaxEvals: 20})
	assert.Equal(t, funcs.MaxEvals, r.Reason)
	assert.GreaterOrEqual(t, r.Evals, 20)
	assert.LessOrEqual(t, r.Evals, 20+2*2)

	r = newBowl(funcs.Stop{GradTol: 1e-3})
	assert.Equal(t, funcs.GradSmall, r.Reason)
	assert.True(t, r.Reason.Converged())
	for _, g := range r.Grad {
		assert.LessOrEqual(t, math.Abs(g), 1e-3)
	}

	// A looser StepTol stops sooner.
	tight := newBowl(funcs.Stop{})
	loose := newBowl(funcs.Stop{StepTol: 1e-3})
	assert.Equal(t, funcs.StepSmall, loose.Reason)
	assert.Less(t, loose.Steps, tight.Steps)

	// Already at the minimum: the gradient is zero before any step.
	d := funcs.NewDescender(funcs.Multi{Ln: 2, M: bowl}, []float64{1, -2})
	d.Stop.GradTol = 1e-9
	r = d.Run()
	assert.Equal(t, funcs.GradSmall, r.Reason)
	assert.Equal(t, 0, r.Steps)

	// With the default Stop, which doesn't check the gradient, it finds that
	// no step lowers f.
	r = funcs.NewDescender(funcs.Multi{Ln: 2, M: bowl}, []float64{1, -2}).Run()
	assert.Equal(t, funcs.StepSmall, r.Reason)
	assert.Equal(t, 0, r.Steps)
	assert.Equal(t, []float64{1, -2}, r.X)

	// The same in one variable, where the numeric gradient is exactly zero.
	var square funcs.M = func(x []float64) float64 { return x[0] * x[0] }
	r = funcs.NewDescender(funcs.Multi{Ln: 1, M: square}, []float64{0}).Run()
	assert.Equal(t, funcs.StepSmall, r.Reason)
	assert.Equal(t, []float64{0}, r.Grad)
}

func TestDescenderNaN(t *testing.T) {
	// f is only defined for x < 3, and its minimum is at 2. The first steps
	// overshoot into the NaN region and the line search backs away from it.
	var f funcs.M = func(x []float64) float64 {
		if x[0] >= 3 {
			return math.NaN()
		}
		return (x[0] - 2) * (x[0] - 2)
	}
	d := funcs.NewDescender(funcs.Multi{Ln: 1, M: f}, []float64{-10})
	d.Direction = &funcs.Gradient{}
	r := d.Run()
	assert.True(t, r.Reason.Converged(), r.Reason.String())
	assert.InDelta(t, 2.0, r.X[0], 1e-6)

	// Starting where f is NaN, there is nowhere to go.
	r = funcs.NewDescender(funcs.Multi{Ln: 1, M: f}, []float64{4}).Run()
	assert.Equal(t, funcs.NaN, r.Reason)
	assert.Equal(t, 0, r.Steps)

	// The same when only the gradient is NaN: √x at 0 is 0, but the numeric
	// derivative needs √-h.
	var sqrt funcs.M = func(x []float64) float64 { return math.Sqrt(x[0]) }
	r = funcs.NewDescender(funcs.Multi{Ln: 1, M: sqrt}, []float64{0}).Run()
	assert.Equal(t, funcs.NaN, r.Reason)
}

// uphill is a Direction that points the wrong way.
type uphill struct{ resets int }

func (u *uphill) Reset(n int) { u.resets++ }
func (u *uphill) Dir(m funcs.M, x []float64, f float64, g, d []float64) float64 {
	copy(d, g)
	return 1
}
func (u *uphill) Update(s, y []float64, t float64) {}

func TestDescenderFallsBackToGradient(t *testing.T) {
	// A direction that doesn't go down is replaced with the gradient, and the
	// Direction is reset.
	u := &uphill{}
	d := funcs.NewDescender(funcs.Multi{Ln: 2, M: bowl}, []float64{0, 0})
	d.Direction = u
	assert.True(t, d.Step())
	assert.Less(t, d.Result().F, 41.0)
	assert.Equal(t, 2, u.resets) // once at the start, once for the fallback
}

func TestLBFGSSkipsNegativeCurvature(t *testing.T) {
	// cos has negative curvature around 0, where the gradient's change shows
	// no positive curvature to learn from. LBFGS skips those steps and still
	// reaches a minimum, an odd multiple of π where cos is -1.
	var f funcs.M = func(x []float64) float64 { return math.Cos(x[0]) }
	r := funcs.NewDescender(funcs.Multi{Ln: 1, M: f}, []float64{0.1}).Run()
	assert.True(t, r.Reason.Converged())
	assert.InDelta(t, -1.0, r.F, 1e-12)
	assert.InDelta(t, 0.0, math.Sin(r.X[0]), 1e-6)
}

func TestReasonString(t *testing.T) {
	for r, s := range map[funcs.Reason]string{
		funcs.Running:   "running",
		funcs.GradSmall: "gradient small",
		funcs.StepSmall: "step small",
		funcs.MaxSteps:  "max steps",
		funcs.MaxEvals:  "max evals",
		funcs.NaN:       "NaN",
		99:              "unknown",
	} {
		assert.Equal(t, s, r.String())
	}
}

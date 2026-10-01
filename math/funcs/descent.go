package funcs

import (
	"math"
)

// == projects.Code.luce.funcs ==
// [ ] Levenberg-Marquardt
//	A solver for systems of equations (residuals that are zero at the
//	solution), such as circuits in melange. It is far less sensitive to
//	badly scaled variables than gradient descent.
// [ ] stochastic descent
//	For functions that are a sum over many terms, step on a sample of the
//	terms at a time.

// https://en.wikipedia.org/wiki/Gradient_descent
// https://en.wikipedia.org/wiki/Limited-memory_BFGS
// https://en.wikipedia.org/wiki/Backtracking_line_search

const (
	// DefaultMaxSteps is the MaxSteps that a zero Stop uses.
	DefaultMaxSteps = 10000
	// DefaultStepTol is the StepTol that a zero Stop uses.
	DefaultStepTol = 1e-15
	// armijo is how much of the decrease that the slope predicts a step has to
	// achieve for the line search to accept it.
	armijo = 1e-4
)

// Stop says when a Descender is done. Zero fields use the defaults.
type Stop struct {
	// MaxSteps limits the number of steps. 0 means DefaultMaxSteps.
	MaxSteps int
	// MaxEvals limits the number of calls to M. 0 means no limit. It is
	// checked before each trial point of the line search, so the gradient at
	// the last accepted point can take it over by up to 2*Ln calls (with
	// NumericDM).
	MaxEvals int
	// GradTol stops when every component of the gradient is within
	// GradTol*max(1, |f|) of zero. 0 means the gradient is not checked: the
	// descent runs until a step can no longer lower f (StepTol). That is the
	// safe choice when the scale of f isn't known, for example a system of
	// equations measured in milliamps.
	GradTol float64
	// StepTol stops when the line search has shrunk the step to within
	// StepTol*max(1, |x|) in every component without lowering f: x is then
	// as close to the minimum as the precision of f allows. 0 means
	// DefaultStepTol.
	StepTol float64
}

func (s Stop) maxSteps() int {
	if s.MaxSteps == 0 {
		return DefaultMaxSteps
	}
	return s.MaxSteps
}

func (s Stop) stepTol() float64 {
	if s.StepTol == 0 {
		return DefaultStepTol
	}
	return s.StepTol
}

// Reason is why a Descender stopped.
type Reason byte

const (
	// Running means the Descender hasn't stopped.
	Running Reason = iota
	// GradSmall means the gradient is within Stop.GradTol of zero.
	GradSmall
	// StepSmall means no step longer than Stop.StepTol lowers f: x is at a
	// minimum, as far as the precision of f allows.
	StepSmall
	// MaxSteps means Stop.MaxSteps steps were taken.
	MaxSteps
	// MaxEvals means Stop.MaxEvals calls to M were made.
	MaxEvals
	// NaN means f or the gradient is NaN at the starting point.
	NaN
)

// Converged is true for the reasons that mean x is at a minimum.
func (r Reason) Converged() bool {
	return r == GradSmall || r == StepSmall
}

func (r Reason) String() string {
	switch r {
	case Running:
		return "running"
	case GradSmall:
		return "gradient small"
	case StepSmall:
		return "step small"
	case MaxSteps:
		return "max steps"
	case MaxEvals:
		return "max evals"
	case NaN:
		return "NaN"
	}
	return "unknown"
}

// Descender minimizes Multi.M, starting from a point, one step at a time. Each
// step asks Direction which way to go, then searches along that direction for
// a point where M is lower by enough (a backtracking line search). A trial
// point where M is NaN counts as too far, so the search backs away from it.
// Stop decides when it is done.
//
// Fields can be set after NewDescender and before the first Step.
type Descender struct {
	Multi
	// Direction picks the direction of each step. nil means
	// &LBFGS{Memory: DefaultMemory}.
	Direction Direction
	Stop      Stop
	// Record, when set, is called with the state before the first step and
	// after each step.
	Record func(StepRecord)

	dm           DM
	x, g         []float64
	d, prevX     []float64
	prevG, s, y  []float64
	f, t         float64
	steps, evals int
	started      bool
	reason       Reason
}

// StepRecord is the state of a Descender after a step. The slices are copies.
type StepRecord struct {
	Step int
	X    []float64
	F    float64
	Grad []float64
	// Dir is the direction of the step that led here, and T the length of the
	// step along it. Both are zero before the first step.
	Dir []float64
	T   float64
}

// Result is the outcome of a descent. X and Grad are the Descender's own
// slices.
type Result struct {
	X, Grad []float64
	F       float64
	Steps   int
	// Evals counts calls to M, including those NumericDM makes.
	Evals  int
	Reason Reason
}

// NewDescender returns a Descender for m that starts at x. x is updated in
// place as the Descender steps.
func NewDescender(m Multi, x []float64) *Descender {
	return &Descender{
		Multi: m,
		x:     x,
	}
}

// start evaluates the starting point. It runs before the first step so that
// fields set after NewDescender are used.
func (d *Descender) start() {
	d.started = true
	if d.Direction == nil {
		d.Direction = &LBFGS{Memory: DefaultMemory}
	}
	n := len(d.x)
	d.Direction.Reset(n)
	counted := Multi{
		Ln: d.Ln,
		M:  d.eval,
		DM: d.Multi.DM,
	}
	d.dm = counted.GetDM()
	d.g = make([]float64, n)
	d.d = make([]float64, n)
	d.prevX = make([]float64, n)
	d.prevG = make([]float64, n)
	d.s = make([]float64, n)
	d.y = make([]float64, n)

	d.f = d.eval(d.x)
	d.g = d.dm(d.x, d.g)
	if math.IsNaN(d.f) || hasNaN(d.g) {
		d.reason = NaN
	} else {
		d.checkGrad()
	}
	d.record()
}

func (d *Descender) eval(x []float64) float64 {
	d.evals++
	return d.Multi.M(x)
}

// Step takes one step. It returns false, without stepping, once the
// Descender has stopped.
func (d *Descender) Step() bool {
	if !d.started {
		d.start()
	}
	if d.reason != Running {
		return false
	}
	if d.steps >= d.Stop.maxSteps() {
		d.reason = MaxSteps
		return false
	}

	d.t = d.Direction.Dir(d.eval, d.x, d.f, d.g, d.d)
	slope := dot(d.g, d.d)
	if !(slope < 0) {
		// Not a way down: fall back to the gradient and start the
		// Direction's history again.
		for i, gi := range d.g {
			d.d[i] = -gi
		}
		slope = dot(d.g, d.d)
		d.t = 1
		d.Direction.Reset(len(d.x))
	}

	copy(d.prevX, d.x)
	copy(d.prevG, d.g)
	tol := d.Stop.stepTol()
	for {
		if d.Stop.MaxEvals > 0 && d.evals >= d.Stop.MaxEvals {
			copy(d.x, d.prevX)
			d.reason = MaxEvals
			return false
		}
		small := true
		for i, di := range d.d {
			d.x[i] = d.prevX[i] + d.t*di
			if math.Abs(d.t*di) > tol*math.Max(1, math.Abs(d.prevX[i])) {
				small = false
			}
		}
		if small {
			copy(d.x, d.prevX)
			d.reason = StepSmall
			return false
		}
		f := d.eval(d.x)
		if f <= d.f+armijo*d.t*slope {
			d.f = f
			break
		}
		d.t /= 2
	}

	d.g = d.dm(d.x, d.g)
	for i := range d.s {
		d.s[i] = d.x[i] - d.prevX[i]
		d.y[i] = d.g[i] - d.prevG[i]
	}
	d.Direction.Update(d.s, d.y, d.t)
	d.steps++
	d.checkGrad()
	d.record()
	return d.reason == Running
}

// checkGrad sets GradSmall if the gradient is within Stop.GradTol of zero.
func (d *Descender) checkGrad() {
	if d.Stop.GradTol == 0 {
		return
	}
	limit := d.Stop.GradTol * math.Max(1, math.Abs(d.f))
	for _, gi := range d.g {
		if math.Abs(gi) > limit {
			return
		}
	}
	d.reason = GradSmall
}

func (d *Descender) record() {
	if d.Record == nil {
		return
	}
	r := StepRecord{
		Step: d.steps,
		X:    append([]float64(nil), d.x...),
		F:    d.f,
		Grad: append([]float64(nil), d.g...),
	}
	if d.steps > 0 {
		r.Dir = append([]float64(nil), d.d...)
		r.T = d.t
	}
	d.Record(r)
}

// Run steps until the Descender stops and returns the Result.
func (d *Descender) Run() Result {
	for d.Step() {
	}
	return d.Result()
}

// Result returns the current state. Its Reason is Running until the
// Descender stops.
func (d *Descender) Result() Result {
	return Result{
		X:      d.x,
		Grad:   d.g,
		F:      d.f,
		Steps:  d.steps,
		Evals:  d.evals,
		Reason: d.reason,
	}
}

func dot(a, b []float64) (sum float64) {
	for i, ai := range a {
		sum += ai * b[i]
	}
	return
}

func hasNaN(xs []float64) bool {
	for _, x := range xs {
		if math.IsNaN(x) {
			return true
		}
	}
	return false
}

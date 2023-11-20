package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
)

const (
	// DefaultDamper is the Damper that Init uses when none is set.
	DefaultDamper = 0.1
)

// == projects.Code.luce.funcs ==
// [ ] gradient descent
//	I need this working for melange
// 	and I'm not sure why it's not
// [ ] stop on NaN
//	Run and Record should stop when encountering NaN

// Note: the issue I was having is that Newton and Secant are for finding zeros.
// I need quasi-Newton methods for finding minima. But that's not necessary
// for Gradient Descent because the single variable will still run the underlying
// function for each step. So it's probably better to just use a straight
// gradient descent.

// variants
// * changing G
// * momentum
// * stochastic

// https://en.wikipedia.org/wiki/Gradient_descent
// https://en.wikipedia.org/wiki/Proximal_gradient_method

// Descender minimizes Multi.M by gradient descent. X is the current position and
// DX the gradient there. Each Step moves X by the momentum, which is the last
// momentum scaled by Damper plus DX scaled by G, then multiplies G by DG,
// recomputes DX and counts Steps down. Init fills in any field left unset.
type Descender struct {
	Multi
	G, DG           float64
	X, DX, Momentum []float64
	Steps           int
	// Momentum is scaled by the damper at each step
	Damper float64
}

// StepRecord is a snapshot of a Descender.
type StepRecord struct {
	G, DG           float64
	X, DX, Momentum []float64
}

// StepRecord returns a copy of the current state. It does not copy Momentum.
func (d *Descender) StepRecord() StepRecord {
	return StepRecord{
		G:  d.G,
		DG: d.DG,
		X:  slice.New(d.X).Clone(-1),
		DX: slice.New(d.DX).Clone(-1),
	}
}

// SetDX sets DX to the gradient at X.
func (d *Descender) SetDX() {
	d.DX = d.Multi.GetDM()(d.X, d.DX)
}

// Step moves X one step and counts Steps down.
func (d *Descender) Step() {
	for i, dx := range d.DX {
		d.Momentum[i] = d.Momentum[i]*d.Damper + dx*d.G
		d.X[i] -= d.Momentum[i]
	}

	d.G *= d.DG
	d.SetDX()
	d.Steps--
}

var absTransform = morph.NewValAll(math.Abs)

// SetG sets G from the largest component m of DX (by absolute value) as the
// smaller of m and max/m.
func (d *Descender) SetG(max float64) *Descender {
	// == projects.Code.luce.funcs ==
	// [ ] Descender.SetG cleaner setup
	//	this is ugly - there should be a cleaner way to set this up
	//	though it is difficult because it can't be a method because the
	//	output type of the transformer is not known
	m := cmpr.MaxN(absTransform.Slice(d.DX, nil)...)
	d.G = math.Min(m, max/m)
	return d
}

// SetDG such that DG^steps = end
func (d *Descender) SetDG(end float64) *Descender {
	d.DG = math.Pow(end, 1.0/float64(d.Steps))
	return d
}

// Init sets fields to reasonable defaults without overriding any values
// that have been set.
func (d *Descender) Init() *Descender {
	if d.X == nil {
		d.X = make([]float64, d.Multi.Ln)
	}
	if d.DX == nil {
		d.DX = make([]float64, d.Multi.Ln)
		d.SetDX()
	}
	if d.Momentum == nil {
		d.Momentum = make([]float64, d.Multi.Ln)
	}
	if d.G == 0 {
		d.SetG(0.1)
	}
	if d.DG == 0 {
		d.SetDG(1e-3)
	}
	if d.Damper == 0 {
		d.Damper = DefaultDamper
	}
	return d
}

// Run steps until Steps is used up.
func (d *Descender) Run() *Descender {
	for d.Steps > 0 {
		d.Step()
	}
	return d
}

// Record is Run that also returns the state before the first step and after
// each step.
func (d *Descender) Record() (*Descender, []StepRecord) {
	log := make([]StepRecord, 0, d.Steps+1)
	log = append(log, d.StepRecord())
	for d.Steps > 0 {
		d.Step()
		log = append(log, d.StepRecord())
	}
	return d, log
}

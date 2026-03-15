package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
)

const (
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

type Descender struct {
	Multi
	G, DG           float64
	X, DX, Momentum []float64
	Steps           int
	// Momentum is scaled by the damper at each step
	Damper float64
}

type StepRecord struct {
	G, DG           float64
	X, DX, Momentum []float64
}

func (d *Descender) StepRecord() StepRecord {
	return StepRecord{
		G:  d.G,
		DG: d.DG,
		X:  slice.New(d.X).Clone(-1),
		DX: slice.New(d.DX).Clone(-1),
	}
}

func (d *Descender) SetDX() {
	d.DX = d.Multi.GetDM()(d.X, d.DX)
}

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

func (d *Descender) SetG(max float64) *Descender {
	// TODO: this is ugly - there should be a cleaner way to set this up
	// though it is difficult because it can't be a method because the
	// output type of the transformer is not known
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

func (d *Descender) Run() *Descender {
	for d.Steps > 0 {
		d.Step()
	}
	return d
}

func (d *Descender) Record() (*Descender, []StepRecord) {
	log := make([]StepRecord, 0, d.Steps+1)
	log = append(log, d.StepRecord())
	for d.Steps > 0 {
		d.Step()
		log = append(log, d.StepRecord())
	}
	return d, log
}

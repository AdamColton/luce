package funcs

import (
	"fmt"
	"math"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
)

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
}

func (d *Descender) SetDX() {
	d.DX = d.Multi.GetDM()(d.X, d.DX)
}

func (d *Descender) Step() {
	//fmt.Print(d.X, d.DX)
	for i, dx := range d.DX {
		d.Momentum[i] = d.Momentum[i]*.1 + dx*d.G
		d.X[i] -= d.Momentum[i]
	}
	fmt.Println(d.Multi.M(d.X))

	d.G *= d.DG
	d.SetDX()
	d.Steps--
}

func (d *Descender) SetG(max float64) *Descender {
	// TODO: this is ugly - there should be a cleaner way to set this up
	// though it is difficult because it can't be a method because the
	// output type of the transformer is not known
	m := cmpr.MaxN(list.NewTransformer(slice.New(d.DX), math.Abs).ToSlice(nil)...)
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
	return d
}

func (d *Descender) Run() *Descender {
	for d.Steps > 0 {
		d.Step()
		fmt.Println(d.Multi.M(d.X))
	}
	return d
}

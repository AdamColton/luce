package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/slice"
)

// CDM: Composable diferential multi func
type CDM interface {
	// M is a multi-func
	M(x []float64) float64

	// IdxDM is the partial derivative for M.
	// The input that it is the derivative of is defined by idx.
	IdxDM(x []float64, idx int) float64
}

type X int

func (i X) M(x []float64) float64 {
	return x[i]
}

func (i X) IdxDM(x []float64, idx int) float64 {
	if i == X(idx) {
		return 1
	}
	return 0
}

type CoExp struct {
	Base CDM
	C, E float64
}

func (ce CoExp) M(x []float64) float64 {
	return ce.C * math.Pow(ce.Base.M(x), ce.E)
}

func (ce CoExp) IdxDM(x []float64, idx int) float64 {
	b := ce.Base.M(x)
	db := ce.Base.IdxDM(x, idx)

	return ce.C * ce.E * db * math.Pow(b, ce.E-1)
}

type Const float64

func (c Const) M([]float64) float64 {
	return float64(c)
}

func (c Const) IdxDM(x []float64, idx int) float64 {
	return 0
}

type Sum []CDM

func (s Sum) M(x []float64) (sum float64) {
	for _, cdm := range s {
		sum += cdm.M(x)
	}
	return
}

func (s Sum) IdxDM(x []float64, idx int) (sum float64) {
	for _, cdm := range s {
		sum += cdm.IdxDM(x, idx)
	}
	return
}

type Product [2]CDM

func (p Product) M(x []float64) float64 {
	return p[0].M(x) * p[1].M(x)
}

func (p Product) IdxDM(x []float64, idx int) float64 {
	p0 := p[0].M(x)
	p1 := p[1].M(x)
	dp0 := p[0].IdxDM(x, idx)
	dp1 := p[1].IdxDM(x, idx)

	return p0*dp1 + p1*dp0
}

type Div [2]CDM

func (d Div) M(x []float64) float64 {
	return d[0].M(x) / d[1].M(x)
}

func (d Div) IdxDM(x []float64, idx int) float64 {
	d0 := d[0].M(x)
	d1 := d[1].M(x)
	dd0 := d[0].IdxDM(x, idx)
	dd1 := d[1].IdxDM(x, idx)

	return (d1*dd0 - d0*dd1) / (d1 * d1)
}

type Ln struct {
	Of CDM
}

func (l Ln) M(x []float64) float64 {
	return math.Log(l.M(x))
}

func (l Ln) IdxDM(x []float64, idx int) float64 {
	m := l.Of.M(x)
	dm := l.Of.IdxDM(x, idx)
	return dm / m
}

type Exp struct {
	Base  float64
	Power CDM
}

func (e Exp) M(x []float64) float64 {
	return math.Pow(e.Base, e.Power.M(x))
}

func (e Exp) IdxDM(x []float64, idx int) float64 {
	m := e.Power.M(x)
	dm := e.Power.IdxDM(x, idx)
	return math.Pow(e.Base, m) * math.Log(e.Base) * dm
}

func RMSE(zeros ...CDM) CDM {
	ln := len(zeros)
	sum := make(Sum, ln)
	for i, z := range zeros {
		sum[i] = CoExp{Base: z, C: 1, E: 2}
	}
	ms := Product{sum, Const(1.0 / float64(ln))}
	return CoExp{Base: ms, C: 1, E: 0.5}
}

func BuildDM(idxDM func(x []float64, idx int) float64) func(x, buf []float64) []float64 {
	return func(x, buf []float64) []float64 {
		buf = slice.NewBuffer(buf).Zeros(len(x))
		for i := range x {
			buf[i] = idxDM(x, i)
		}
		return buf
	}
}

// System expects the zeros CDMs to be equal to zero when a correct solution
// is found. It wraps the whole system of equations in a RMSE.
func System(x []float64, zeros []CDM) *Descender {
	ln := len(x)
	rmse := RMSE(zeros...)
	return (&Descender{
		Multi: Multi{
			Ln: ln,
			M:  rmse.M,
			DM: BuildDM(rmse.IdxDM),
		},
		Steps: 1000,
		X:     x,
	})
}

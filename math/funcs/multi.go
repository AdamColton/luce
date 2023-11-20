package funcs

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
)

type M func([]float64) float64

func (fn M) PartialDerivative(x []float64, idx int) float64 {
	// d is the delta between x+step and x-step
	// first we're looking for a step size that gets d close to zero
	const small cmpr.Tolerance = 1e-3
	step := 1e-2
	d := 1.0
	xi := x[idx]
	for !small.Zero(d) {
		step /= 2
		x[idx] = xi + step
		d = fn(x)
		x[idx] = xi - step
		d -= fn(x)
		x[idx] = xi
	}
	return d / (2 * step)
}

type DM func(x, buf []float64) []float64

type Multi struct {
	Ln int
	M  M
	DM DM
}

func (s *Multi) AnalyticDM(x, buf []float64) []float64 {
	out := slice.NewBuffer(buf).Slice(s.Ln)
	for i := range out {
		out[i] = s.M.PartialDerivative(x, i)
	}

	return out
}

func (s *Multi) GetDM() DM {
	if s.DM == nil {
		return s.AnalyticDM
	}
	return s.DM
}

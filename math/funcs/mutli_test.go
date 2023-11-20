package funcs_test

import (
	"math"
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/math/funcs"
	"github.com/adamcolton/luce/math/ints/comb"
	"github.com/stretchr/testify/assert"
)

var (
	circA = func(t float64) (x, y float64) {
		x, y = math.Sincos(t * 6.28)
		return
	}
	circB = func(t float64) (x, y float64) {
		x, y = math.Sincos(t*6.28 + 1)
		x += 0.5
		return
	}
	dist funcs.M = func(t []float64) float64 {
		x0, y0 := circA(t[0])
		x1, y1 := circB(t[1])
		dx := (x0 - x1)
		dy := (y0 - y1)
		d := math.Sqrt(dx*dx + dy*dy)
		return d
	}
	dcircA = func(t float64) (x, y float64) {
		y, x = math.Sincos(t * 6.28)
		y = -y * 6.28
		x = x * 6.28
		return
	}
	dcircB = func(t float64) (x, y float64) {
		y, x = math.Sincos(t*6.28 + 1)
		y = -y * 6.28
		x = x * 6.28
		return
	}
	ddist = func(t0, t1 float64) (float64, float64) {
		// d(t0,t1) = ( (ax-bx)^2 + (ay-by)^2 )^0.5 [where ax,ay = ca(t0); bx,by = cb(t1)]
		// d = f(g)
		// f = g^0.5 :. f'(g) = ((g^-0.5)/2) * g'
		// g(h,i) = h^2 + i^2 :. g'(h,i) = 2h*h' + 2i*i'
		// h(t0,t1) = ca(t0).x-cb(t1).x :. dh0 = dca(t0).x; dh1 = -dcb(t1).x
		// i(t0,t1) =  ca(t0).y-cb(t1).y :. di0 = dca(t0).y; di1 = -dcb(t1).y
		// dg0 = 2h*dh0 + 2i*di0 ; dg1 = 2h*dh1 + 2i*di1
		// gc = ((g^-0.5)/2)
		// df0 = gc * dg0 ; df1 = gc*dg1

		ax, ay := circA(t0)
		bx, by := circB(t1)
		h := ax - bx
		i := ay - by
		g := h*h + i*i

		dh0, di0 := dcircA(t0)
		dh1, di1 := dcircB(t1)
		dh1, di1 = -dh1, -di1

		dg0 := 2*h*dh0 + 2*i*di0
		dg1 := 2*h*dh1 + 2*i*di1
		gc := math.Pow(g, -0.5) / 2

		df0 := gc * dg0
		df1 := gc * dg1

		return df0, df1
	}
	resolution = 10.0
	steps      = list.Generator[float64]{
		Length: int(resolution),
		Fn: func(i int) float64 {
			return float64(i) / resolution
		},
	}
)

func TestSanityCheck(t *testing.T) {
	// This test just confirms that the reference functions are correct
	d := 1e-5
	check := func(fn, dfn func(t float64) (x, y float64)) {
		for i := range steps.Wrap().Iter().Channel(0) {
			x, y := fn(i)
			dx, dy := dfn(i)
			x0, y0 := fn(i + d)

			assert.InDelta(t, dx, (x0-x)/d, 1e-3)
			assert.InDelta(t, dy, (y0-y)/d, 1e-3)
		}
	}

	check(circA, dcircA)
	check(circB, dcircB)

	grid := list.Combinator(steps, steps, comb.Cross)
	for pt := range grid.Iter().Channel(0) {
		y := dist([]float64{pt.A, pt.B})
		yA := dist([]float64{pt.A + d, pt.B})
		yB := dist([]float64{pt.A, pt.B + d})
		dya, dyb := ddist(pt.A, pt.B)

		assert.InDelta(t, dya, (yA-y)/d, 5e-3)
		assert.InDelta(t, dyb, (yB-y)/d, 5e-3)
	}
}

func TestPartialDerivative(t *testing.T) {
	m := funcs.Multi{
		Ln: 2,
		M:  dist,
	}

	p := make([]float64, 2)
	grid := list.NewTransformer(
		list.Combinator(steps, steps, comb.Cross),
		func(s struct{ A, B float64 }) []float64 {
			p[0], p[1] = s.A, s.B
			return p
		},
	)
	buf := make([]float64, 2)
	expected := make([]float64, 2)
	dm := m.GetDM()
	for pt := range grid.Iter().Channel(0) {
		expected[0], expected[1] = ddist(pt[0], pt[1])
		got := dm(pt, buf)
		t.Log(pt)
		for i, g := range got {
			assert.InDelta(t, expected[i], g, 4e-4)
		}
	}
}

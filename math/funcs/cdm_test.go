package funcs_test

import (
	"math"
	"testing"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/stretchr/testify/assert"
)

func TestCDM(t *testing.T) {
	var cdm funcs.CDM = funcs.CoExp{
		Base: funcs.X(0),
		C:    1.5,
		E:    2,
	}

	cdm = funcs.Sum{
		cdm,
		funcs.Const(4),
	}
	x := cdm.M([]float64{3})
	assert.Equal(t, 17.5, x)

	x = cdm.IdxDM([]float64{3}, 0)
	assert.Equal(t, 9.0, x)
	assert.Equal(t, "1.5*x[0]^2 + 4", cdm.(funcs.Sum).String())
}

// checkCDM checks that c has the value want at each x, and that its exact
// partial derivatives agree with numeric ones there.
func checkCDM(t *testing.T, c funcs.CDM, want func(x []float64) float64, xs ...[]float64) {
	t.Helper()
	m := funcs.Multi{Ln: len(xs[0]), M: c.M, DM: funcs.BuildDM(c.IdxDM)}
	for _, x := range xs {
		assert.InDelta(t, want(x), c.M(x), 1e-12*math.Max(1, math.Abs(want(x))), x)
		ok, idx := m.DMAgrees(x, funcs.DefaultCheckTol)
		assert.True(t, ok, "partial %d at %v", idx, x)
	}
}

func TestCDMs(t *testing.T) {
	x0, x1, x2 := funcs.X(0), funcs.X(1), funcs.X(2)
	tests := map[string]struct {
		c    funcs.CDM
		want func(x []float64) float64
		str  string
	}{
		"diode": {
			// The Shockley diode equation: the current through a diode for a
			// voltage x0, with a saturation current of 1e-12 A and a thermal
			// voltage of 25 mV.
			c: funcs.Scale(1e-12, funcs.Sub(
				funcs.Exp(funcs.Scale(1/0.025, x0)),
				funcs.Const(1),
			)),
			want: func(x []float64) float64 {
				return 1e-12 * (math.Exp(x[0]/0.025) - 1)
			},
			str: "1e-12*(exp(40*x[0]) - 1)",
		},
		"volume": {
			c: funcs.Product{x0, x1, x2},
			want: func(x []float64) float64 {
				return x[0] * x[1] * x[2]
			},
			str: "x[0]*x[1]*x[2]",
		},
		"scaled-square": {
			// Only the factor x1² depends on x1; the rest are constant.
			c: funcs.Product{funcs.Const(3), funcs.Square(x1), x0},
			want: func(x []float64) float64 {
				return 3 * x[1] * x[1] * x[0]
			},
			str: "3*(x[1]^2)*x[0]",
		},
		"damped": {
			// x0 appears in both factors.
			c: funcs.Product{x0, funcs.Exp(funcs.Neg(x0)), x1},
			want: func(x []float64) float64 {
				return x[0] * math.Exp(-x[0]) * x[1]
			},
			str: "x[0]*exp(-x[0])*x[1]",
		},
		"parallel": {
			// Two resistors in parallel: x0*x1/(x0 + x1).
			c: funcs.Div{funcs.Product{x0, x1}, funcs.Add(x0, x1)},
			want: func(x []float64) float64 {
				return x[0] * x[1] / (x[0] + x[1])
			},
			str: "(x[0]*x[1])/(x[0] + x[1])",
		},
		"per-unit": {
			c: funcs.Div{x0, funcs.Const(4)},
			want: func(x []float64) float64 {
				return x[0] / 4
			},
			str: "x[0]/4",
		},
		"log": {
			c: funcs.Ln{Of: funcs.Add(funcs.Square(x0), funcs.Const(1))},
			want: func(x []float64) float64 {
				return math.Log(x[0]*x[0] + 1)
			},
			str: "ln(x[0]^2 + 1)",
		},
		"node": {
			// The current into a node from 5 V through 100 Ω, and out through
			// 200 Ω, with x0 the node voltage and x1 a current source.
			c: funcs.Linear{Coef: []float64{-1.0/100 - 1.0/200, 1}, C: 5.0 / 100},
			want: func(x []float64) float64 {
				return 0.05 - 0.015*x[0] + x[1]
			},
			str: "-0.015*x[0] + x[1] + 0.05",
		},
		"wave": {
			c: funcs.Add(funcs.Sin(x0), funcs.Cos(funcs.Scale(2, x1))),
			want: func(x []float64) float64 {
				return math.Sin(x[0]) + math.Cos(2*x[1])
			},
			str: "sin(x[0]) + cos(2*x[1])",
		},
		"bounded": {
			// x0 kept between 1 and 3, and x1 between -1 and 1.
			c: funcs.Add(
				funcs.Add(funcs.Const(1), funcs.Scale(2, funcs.Sigmoid(x0))),
				funcs.Tanh(x1),
			),
			want: func(x []float64) float64 {
				return 1 + 2/(1+math.Exp(-x[0])) + math.Tanh(x[1])
			},
			str: "1 + 2*sigmoid(x[0]) + tanh(x[1])",
		},
		"distance": {
			c: funcs.Sqrt(funcs.Add(funcs.Square(x0), funcs.Square(x1))),
			want: func(x []float64) float64 {
				return math.Hypot(x[0], x[1])
			},
			str: "sqrt(x[0]^2 + x[1]^2)",
		},
		"abs": {
			c: funcs.Abs(funcs.Sub(x0, x1)),
			want: func(x []float64) float64 {
				return math.Abs(x[0] - x[1])
			},
			str: "abs(x[0] - x[1])",
		},
		"rmse": {
			c: funcs.RMSE(funcs.Sub(x0, funcs.Const(1)), funcs.Sub(x1, funcs.Const(2))),
			want: func(x []float64) float64 {
				a, b := x[0]-1, x[1]-2
				return math.Sqrt((a*a + b*b) / 2)
			},
			str: "sqrt(((x[0] - 1)^2 + (x[1] - 2)^2)*0.5)",
		},
		"negative": {
			c: funcs.Neg(funcs.Product{x0, funcs.Const(2)}),
			want: func(x []float64) float64 {
				return -2 * x[0]
			},
			str: "-(x[0]*2)",
		},
	}

	// Points away from where any of these isn't differentiable or defined: the
	// variables are positive and none are equal, nor equal to 0 or each other.
	xs := [][]float64{
		{0.3, 1.7, 2.2},
		{1.4, 0.6, 0.9},
		{2.5, 3.1, 0.4},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkCDM(t, tc.c, tc.want, xs...)
			assert.Equal(t, tc.str, tc.c.(interface{ String() string }).String())
		})
	}
}

func TestLn(t *testing.T) {
	// ln(x0²) at x0 = 3, and its derivative 2/x0.
	ln := funcs.Ln{Of: funcs.CoExp{Base: funcs.X(0), C: 1, E: 2}}
	x := []float64{3}
	assert.InDelta(t, math.Log(9), ln.M(x), 1e-12)
	assert.InDelta(t, 2.0/3, ln.IdxDM(x, 0), 1e-12)
}

// counted is the variable x[idx], as a CDM that counts calls to M, to show
// which ones are skipped. It has no String method.
type counted struct {
	idx   int
	calls *int
}

func (c counted) M(x []float64) float64 {
	*c.calls++
	return x[c.idx]
}

func (c counted) IdxDM(x []float64, idx int) float64 {
	if idx == c.idx {
		return 1
	}
	return 0
}

func TestCDMSkipsUnusedValues(t *testing.T) {
	// When a term doesn't depend on x[idx], its partial derivative is 0 and
	// its value isn't needed. In a circuit, most equations use only a few of
	// the variables, so most partial derivatives are found without
	// evaluating anything.
	calls := 0
	x0 := counted{idx: 0, calls: &calls}
	x := []float64{3, 5}
	cdms := []funcs.CDM{
		funcs.Product{x0, x0},
		funcs.Div{x0, x0},
		funcs.Ln{Of: x0},
		funcs.Square(x0),
		funcs.Sin(x0),
	}
	for _, c := range cdms {
		assert.Equal(t, 0.0, c.IdxDM(x, 1))
	}
	assert.Equal(t, 0, calls)
}

func TestLinear(t *testing.T) {
	// Variables past the end of Coef have a coefficient of 0.
	l := funcs.Linear{Coef: []float64{2, 0, -1}}
	x := []float64{1, 2, 3, 4}
	assert.Equal(t, -1.0, l.M(x))
	assert.Equal(t, 0.0, l.IdxDM(x, 3))
	assert.Equal(t, "2*x[0] - x[2]", l.String())
	assert.Equal(t, "0", funcs.Linear{}.String())
}

func TestCDMStringWithoutStringer(t *testing.T) {
	// A CDM without a String method shows as "?".
	c := funcs.Add(funcs.X(0), counted{})
	assert.Equal(t, "x[0] + ?", c.(funcs.Sum).String())
}

func TestSystem(t *testing.T) {
	// A circuit: 5V through 100Ω to node n, then 200Ω and 300Ω in parallel to
	// ground. x[0] is the current from the source and x[1] the voltage at n.
	// Each equation is the current into a node, which is zero when solved.
	node1 := funcs.Sum{
		funcs.X(0),
		funcs.CoExp{Base: funcs.X(1), C: 1.0 / 100, E: 1},
		funcs.Const(-5.0 / 100),
	}
	node2 := funcs.Sum{
		funcs.Const(5.0 / 100),
		funcs.CoExp{Base: funcs.X(1), C: -1.0 / 100, E: 1},
		funcs.CoExp{Base: funcs.X(1), C: -1.0 / 200, E: 1},
		funcs.CoExp{Base: funcs.X(1), C: -1.0 / 300, E: 1},
	}

	r := funcs.System([]float64{0, 0}, []funcs.CDM{node1, node2}).Run()
	assert.True(t, r.Reason.Converged(), r.Reason.String())
	// 200Ω and 300Ω in parallel are 120Ω, so 5V across 220Ω in all.
	assert.InDelta(t, 5.0/220, r.X[0], 1e-9)
	assert.InDelta(t, 5*120.0/220, r.X[1], 1e-9)

	// SumSquares is what System minimizes: half the sum of the squares.
	ss := funcs.SumSquares(node1, node2)
	x := []float64{0, 0}
	assert.InDelta(t, (0.05*0.05+0.05*0.05)/2, ss.M(x), 1e-15)
}

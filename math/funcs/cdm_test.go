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

}

func TestLn(t *testing.T) {
	// ln(x0²) at x0 = 3, and its derivative 2/x0.
	ln := funcs.Ln{Of: funcs.CoExp{Base: funcs.X(0), C: 1, E: 2}}
	x := []float64{3}
	assert.InDelta(t, math.Log(9), ln.M(x), 1e-12)
	assert.InDelta(t, 2.0/3, ln.IdxDM(x, 0), 1e-12)
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

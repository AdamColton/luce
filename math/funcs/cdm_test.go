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

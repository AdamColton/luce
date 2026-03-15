package funcs_test

import (
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

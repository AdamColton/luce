package funcs_test

import (
	"testing"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/stretchr/testify/assert"
)

func TestGradientDescent(t *testing.T) {
	d := (&funcs.Descender{
		Multi: funcs.Multi{
			Ln: 2,
			M:  dist,
		},
		Steps: 25,
		X:     []float64{0.4, 0.3},
	}).
		Init().
		Run()
	assert.True(t, d.M(d.X) < 1e-3)
}

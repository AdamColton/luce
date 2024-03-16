package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/math/cmpr/cmprtest"
	"github.com/stretchr/testify/assert"
)

func TestAssertEqual(t *testing.T) {
	tt := map[string]struct {
		A, B  any
		equal bool
	}{
		"wrapper-wrapper-equal": {
			A:     list.Slice([]float64{3, 1, 4, 1, 5}),
			B:     list.Slice([]float64{3, 1, 4, 1, 5}),
			equal: true,
		},
		"wrapper-wrapper-not-equal": {
			A:     list.Slice([]float64{3, 1, 4, 1, 5}),
			B:     list.Slice([]float64{3, 1, 5, 1, 5}),
			equal: false,
		},
		"wrapper-wrapper-too-short": {
			A:     list.Slice([]float64{3, 1, 4, 1, 5}),
			B:     list.Slice([]float64{3, 1, 4, 1}),
			equal: false,
		},
		"wrapper-wrapper-too-long": {
			A:     list.Slice([]float64{3, 1, 4, 1, 5}),
			B:     list.Slice([]float64{3, 1, 4, 1, 5, 9}),
			equal: false,
		},
		"wrapper-wrapper-slice-equal": {
			A:     list.Slice([]float64{3, 1, 4, 1, 5}),
			B:     []float64{3.0000001, 1, 4, 1, 5},
			equal: true,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			err := cmprtest.AssertEqual(tc.A, tc.B, 1e-6)
			assert.Equal(t, tc.equal, err == nil)
		})
	}
}

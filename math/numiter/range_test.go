package numiter_test

import (
	"testing"

	"github.com/adamcolton/luce/math/numiter"
	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	tt := map[string]struct {
		expected []float64
		r        *numiter.Range[float64]
	}{
		"NewRange": {
			expected: []float64{0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5},
			r:        numiter.NewRange(0.0, 4.0, 0.5),
		},
		"Include": {
			expected: []float64{0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4.0},
			r:        numiter.Include(0.0, 3.8, 0.5),
		},
		"IntRange": {
			expected: []float64{0, 1, 2, 3},
			r:        numiter.IntRange(4.0),
		},
		"Float(1/3)|Regression": {
			expected: []float64{1.0 / 6.0, 3.0 / 6.0, 5.0 / 6.0},
			r:        numiter.NewRange(1.0/6.0, 1, 1.0/3.0),
		},
		"Float(1/3)+|Regression": {
			expected: []float64{3.0 / 12.0, 7.0 / 12.0, 11.0 / 12.0},
			r:        numiter.NewRange(3.0/12.0, 1, 1.0/3.0),
		},
		"Float(1/3)-|Regression": {
			expected: []float64{1.0 / 12.0, 5.0 / 12.0, 9.0 / 12.0},
			r:        numiter.NewRange(1.0/12.0, 1, 1.0/3.0),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			ln := tc.r.Len()
			if assert.Equal(t, len(tc.expected), ln) {
				for i := range ln {
					assert.InDelta(t, tc.expected[i], tc.r.AtIdx(i), 1e-8)
				}
			}
		})
	}

	assert.Equal(t, 3, numiter.NewRange(0, 3, 1).Len())
	assert.Equal(t, 2, numiter.NewRange(0, 4, 2).Len())
}

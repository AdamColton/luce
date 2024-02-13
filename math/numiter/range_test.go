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
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			ln := tc.r.Len()
			if assert.Equal(t, len(tc.expected), ln) {
				for i := 0; i < ln; i++ {
					assert.Equal(t, tc.expected[i], tc.r.AtIdx(i))
				}
			}
		})
	}
}

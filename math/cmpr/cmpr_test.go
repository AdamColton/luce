package cmpr_test

import (
	"testing"

	"github.com/adamcolton/luce/math/cmpr"
	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	assert.Equal(t, 5, cmpr.Min(5, 6))
	assert.Equal(t, 6, cmpr.Min(7, 6))
	assert.Equal(t, "good-bye", cmpr.Min("hello", "good-bye"))
}

func TestMinN(t *testing.T) {
	assert.Equal(t, 1, cmpr.MinN(3, 1, 4, 1, 5, 9))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 6, cmpr.Max(5, 6))
	assert.Equal(t, 7, cmpr.Max(7, 6))
	assert.Equal(t, "hello", cmpr.Max("hello", "good-bye"))
}

func TestMaxN(t *testing.T) {
	assert.Equal(t, 9, cmpr.MaxN(3, 1, 4, 1, 5, 9))
	assert.Equal(t, 0, cmpr.MaxN[int]())
}

func TestUnique(t *testing.T) {
	tt := map[string]struct {
		in, expected []float64
		cmpr.Tolerance
	}{
		"no-change": {
			in:        []float64{1, 2, 3, 4, 5},
			expected:  []float64{1, 2, 3, 4, 5},
			Tolerance: 1e-3,
		},
		"tail-equal": {
			in:        []float64{1, 2, 3, 4, 5, 5},
			expected:  []float64{1, 2, 3, 4, 5},
			Tolerance: 1e-3,
		},
		"head-equal": {
			in:        []float64{1, 1, 2, 3, 4, 5},
			expected:  []float64{1, 2, 3, 4, 5},
			Tolerance: 1e-3,
		},
		"double": {
			in:        []float64{1, 1, 2, 2, 3, 3, 4, 4, 5, 5},
			expected:  []float64{1, 2, 3, 4, 5},
			Tolerance: 1e-3,
		},
		"close": {
			in:        []float64{1, 2, 3, 3.0006, 4, 5},
			expected:  []float64{1, 2, 3, 4, 5},
			Tolerance: 1e-3,
		},
		"close-chain": {
			in:        []float64{1, 2, 3, 3.0006, 3.0012, 3.0018, 4, 5},
			expected:  []float64{1, 2, 3, 3.0012, 4, 5},
			Tolerance: 1e-3,
		},
		"zero": {
			in:        []float64{},
			expected:  []float64{},
			Tolerance: 1e-3,
		},
		"one": {
			in:        []float64{1},
			expected:  []float64{1},
			Tolerance: 1e-3,
		},
		"regression": {
			in:        []float64{0, 0, 1.0 / 3.0},
			expected:  []float64{0, 1.0 / 3.0},
			Tolerance: 1e-6,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			got := tc.Unique(tc.in)
			assert.Equal(t, tc.expected, got)
		})
	}
}

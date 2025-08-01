package numiter_test

import (
	"testing"

	"github.com/adamcolton/geom/angle"
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr/cmprtest"
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
		"Steps": {
			expected: []float64{0, 0.2, 0.4, 0.6, 0.8},
			r:        numiter.Steps(0.0, 1.0, 5, false),
		},
		"StepsIncludeEnd": {
			expected: []float64{0, 0.25, 0.5, 0.75, 1.0},
			r:        numiter.Steps(0.0, 1.0, 5, true),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			cmprtest.Equal(t, list.Slice(tc.expected), tc.r)
		})
	}

	assert.Equal(t, 3, numiter.NewRange(0, 3, 1).Len())
	assert.Equal(t, 2, numiter.NewRange(0, 4, 2).Len())
}

func TestBadGrid(t *testing.T) {
	defer func() {
		assert.Equal(t, numiter.ErrBadGrid, recover())
	}()

	numiter.Grid(1, 2)
}

func TestGrid(t *testing.T) {
	expected := [][]float64{
		{0, 0}, {0.5, 0}, {1, 0}, {1.5, 0},
		{0, 0.25}, {0.5, 0.25}, {1, 0.25}, {1.5, 0.25},
		{0, 0.5}, {0.5, 0.5}, {1, 0.5}, {1.5, 0.5},
		{0, 0.75}, {0.5, 0.75}, {1, 0.75}, {1.5, 0.75},
		{0, 1}, {0.5, 1}, {1, 1}, {1.5, 1},
		{0, 1.25}, {0.5, 1.25}, {1, 1.25}, {1.5, 1.25},
		{0, 1.5}, {0.5, 1.5}, {1, 1.5}, {1.5, 1.5},
		{0, 1.75}, {0.5, 1.75}, {1, 1.75}, {1.5, 1.75},
		{0, 2}, {0.5, 2}, {1, 2}, {1.5, 2},
		{0, 2.25}, {0.5, 2.25}, {1, 2.25}, {1.5, 2.25},
		{0, 2.5}, {0.5, 2.5}, {1, 2.5}, {1.5, 2.5},
		{0, 2.75}, {0.5, 2.75}, {1, 2.75}, {1.5, 2.75},
	}
	cmprtest.Equal(t, numiter.Grid(0, 2, .5, 0, 3, .25), expected)
}

func TestIntGrid(t *testing.T) {
	expected := [][]int{
		{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {1, 1, 0}, {0, 2, 0}, {1, 2, 0},
		{0, 0, 1}, {1, 0, 1}, {0, 1, 1}, {1, 1, 1}, {0, 2, 1}, {1, 2, 1},
		{0, 0, 2}, {1, 0, 2}, {0, 1, 2}, {1, 1, 2}, {0, 2, 2}, {1, 2, 2},
		{0, 0, 3}, {1, 0, 3}, {0, 1, 3}, {1, 1, 3}, {0, 2, 3}, {1, 2, 3},
	}
	cmprtest.Equal(t, numiter.IntGrid(2, 3, 4), expected)
}

func TestSteps(t *testing.T) {
	rng := numiter.Steps(0, angle.Rot(1), 100, false)
	assert.Equal(t, 100, rng.Len())
	i := rng.Iter()
	c := 0
	for _, done := i.Cur(); !done; _, done = i.Next() {
		c++
	}
	assert.Equal(t, 100, c)
}

func TestIndex(t *testing.T) {
	i := slice.NewIndex(5, 3)
	r := numiter.Index(i)
	expected := &numiter.Range[int]{
		Start: 5,
		End:   8,
		Step:  1,
	}
	assert.Equal(t, expected, r)
}

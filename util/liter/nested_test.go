package liter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestNested(t *testing.T) {
	s := &sliceIter[int]{
		Slice: []int{-1, 1, 2, 3, -1, 4, 5, 6},
	}
	lkup := func(i int) liter.Iter[float64] {
		if i == -1 {
			return &sliceIter[float64]{
				Slice: []float64{},
			}
		}
		f := float64(i)
		return &sliceIter[float64]{
			Slice: []float64{f * 0.5, f, f * 1.5, f * 2},
		}
	}
	n := liter.NewNested(s, lkup).Wrap()
	assert.False(t, n.Done())

	got := make([]float64, 0, 6*4)
	n.For(func(f float64) {
		assert.Equal(t, len(got), n.Idx())
		got = append(got, f)
		gotf, _ := n.Cur()
		assert.Equal(t, f, gotf)
	})

	expected := []float64{
		0.5, 1, 1.5, 2,
		1.0, 2, 3.0, 4,
		1.5, 3, 4.5, 6,
		2.0, 4, 6.0, 8,
		2.5, 5, 7.5, 10,
		3.0, 6, 9.0, 12,
	}
	assert.Equal(t, expected, got)

	_, done := n.Cur()
	assert.True(t, done)

	s = &sliceIter[int]{
		Slice: []int{},
	}
	n = liter.NewNested(s, lkup).Wrap()
	_, done = n.Cur()
	assert.True(t, done)
}

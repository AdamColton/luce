package liter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func (s *sliceIter[T]) Wrap() liter.Wrapper[T] {
	return liter.Wrapper[T]{s}
}

func TestWrap(t *testing.T) {
	si := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	w := liter.Wrap[int](si)
	assert.Equal(t, si, w.Iter)
	w = liter.Wrap[int](w)
	assert.Equal(t, si, w.Iter)
}

package iter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/iter"
	"github.com/stretchr/testify/assert"
)

func (s *sliceIter[T]) Wrap() iter.Wrapper[T] {
	return iter.Wrapper[T]{s}
}

func TestWrap(t *testing.T) {
	si := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	w := iter.Wrap[int](si)
	assert.Equal(t, si, w.Iter)
	w = iter.Wrap[int](w)
	assert.Equal(t, si, w.Iter)
}

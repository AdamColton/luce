package iter_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

func (s *sliceIter[T]) Wrap() iter.Wrapper[T] {
	return iter.Wrapper[T]{s}
}

func (s *sliceIter[T]) String() string {
	return "sliceIter"
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

func TestUpgrade(t *testing.T) {
	si := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	w := si.Wrap()

	s, ok := upgrade.To[fmt.Stringer](w)
	assert.True(t, ok)
	assert.NotNil(t, s)
	assert.Equal(t, "sliceIter", s.String())
}

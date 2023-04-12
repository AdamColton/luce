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

	var s fmt.Stringer
	upgrade.Upgrade(w, &s)
	assert.NotNil(t, s)
	assert.Equal(t, "sliceIter", s.String())
}

func TestWrapperSeek(t *testing.T) {
	si := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	w := si.Wrap()
	idx := 0
	fn := func(i int) bool {
		assert.Equal(t, si.Slice[idx], i)
		idx++
		return false
	}
	it := w.Seek(fn)
	assert.Len(t, si.Slice, idx)
	assert.Nil(t, it)

	w.Iter.(*sliceIter[int]).idx = 0
	it = w.Seek(func(i int) bool {
		return i == 4
	})
	i, done := it.Cur()
	assert.Equal(t, 4, i)
	assert.False(t, done)
	idx = si.idx

	it = w.Seek(fn)
	assert.Len(t, si.Slice, idx)
	assert.Nil(t, it)
}

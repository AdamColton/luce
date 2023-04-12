package liter_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

func (s *sliceIter[T]) Wrap() liter.Wrapper[T] {
	return liter.Wrapper[T]{s}
}

func (s *sliceIter[T]) String() string {
	return "sliceIter"
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

func TestWrapperFor(t *testing.T) {
	si := &sliceIter[byte]{
		Slice: []byte("hello"),
	}
	w := si.Wrap()
	out := ""
	fn := func(b byte) {
		out += string(b)
	}
	w.For(fn)
	assert.Equal(t, "hello", out)
}

func TestWrapperForIdx(t *testing.T) {
	si := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	w := si.Wrap()
	fn := func(i, idx int) {
		assert.Equal(t, si.Slice[idx], i)
	}
	c := w.ForIdx(fn)
	assert.Len(t, si.Slice, c)
}

func TestWrapperConcurrent(t *testing.T) {
	si := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	w := si.Wrap()
	var c int32
	wg := w.Concurrent(func(i, idx int) {
		assert.Equal(t, si.Slice[idx], i)
		atomic.AddInt32(&c, 1)
	})
	wg.Wait()
	assert.Len(t, si.Slice, int(c))
}

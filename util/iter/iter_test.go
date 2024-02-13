package iter_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/adamcolton/luce/util/iter"
	"github.com/stretchr/testify/assert"
)

type sliceIter[T any] struct {
	Slice []T
	idx   int
}

func (s *sliceIter[T]) Next() (T, bool) {
	s.idx++
	return s.Cur()
}
func (s *sliceIter[T]) Cur() (t T, done bool) {
	done = s.
		Done()
	if !done {
		t = s.Slice[s.idx]
	}
	return
}
func (s *sliceIter[T]) Done() bool {
	return s.idx >= len(s.Slice)
}
func (s *sliceIter[T]) Idx() int {
	return s.idx
}

func TestIterSeek(t *testing.T) {
	s := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	idx := 0
	fn := func(i int) bool {
		assert.Equal(t, s.Slice[idx], i)
		idx++
		return false
	}
	it := iter.Seek[int](s, fn)
	assert.Len(t, s.Slice, idx)
	assert.Nil(t, it)

	s.idx = 0
	it = iter.Seek[int](s, func(i int) bool {
		return i == 4
	})
	i, done := it.Cur()
	assert.Equal(t, 4, i)
	assert.False(t, done)
	idx = s.idx

	it = iter.Seek[int](s, fn)
	assert.Len(t, s.Slice, idx)
	assert.Nil(t, it)
}

func ExampleSeek() {
	s := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	fn := func(i int) bool {
		return i == 4
	}
	it := iter.Seek[int](s, fn)

	v, _ := it.Cur()
	fmt.Println("Value:", v, "Idx:", it.Idx())

	// Output:
	// Value: 4 Idx: 2
}

func TestIterFor(t *testing.T) {
	expected := "hello"
	si := &sliceIter[byte]{
		Slice: []byte(expected),
	}
	out := ""
	fn := func(b byte) {
		out += string(b)
	}
	iter.For[byte](si, fn)
	assert.Equal(t, expected, out)
}

func ExampleFor() {
	si := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	fn := func(i int) {
		fmt.Println(i)
	}
	iter.For[int](si, fn)
	// Output:
	// 3
	// 1
	// 4
	// 1
	// 5
	// 9
}

func TestIterForIdx(t *testing.T) {
	s := []int{3, 1, 4, 1, 5, 9}
	si := &sliceIter[int]{
		Slice: s,
	}
	fn := func(i, idx int) {
		assert.Equal(t, s[idx], i)
	}
	c := iter.ForIdx[int](si, fn)
	assert.Len(t, s, c)
}

func TestIterConcurrent(t *testing.T) {
	s := &sliceIter[int]{
		Slice: make([]int, 1000),
	}
	for i := range s.Slice {
		s.Slice[i] = i
	}
	c := atomic.Int32{}
	fn := func(i, idx int) {
		assert.Equal(t, s.Slice[idx], i)
		c.Add(1)
	}
	iter.Concurrent(s, fn).Wait()
	assert.Len(t, s.Slice, int(c.Load()))

	c.Store(0)
	iter.Concurrent(s, fn).Wait()
	assert.Equal(t, 0, int(c.Load()))

	s.idx = 0
	iter.Concurrent(s, func(i, idx int) {
		s.Slice[idx] -= i
	}).Wait()
	for _, i := range s.Slice {
		assert.Equal(t, 0, i)
	}
}

func TestIterChannel(t *testing.T) {
	s := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}

	idx := 0
	for i := range iter.Channel[int](s, 0) {
		assert.Equal(t, s.Slice[idx], i)
		idx++
	}
	assert.Len(t, s.Slice, idx)
}

func sliceFactory[T any](slice []T) iter.Factory[T] {
	return func() (i iter.Iter[T], t T, done bool) {
		i = &sliceIter[T]{
			Slice: slice,
		}
		t, done = i.Cur()
		return
	}
}

func TestFactorySeek(t *testing.T) {
	s := []int{3, 1, 4, 1, 5, 9}
	var sf iter.Factory[int] = sliceFactory(s)
	idx := 0
	fn := func(i int) bool {
		assert.Equal(t, s[idx], i)
		idx++
		return false
	}
	it := sf.Seek(fn)
	assert.Len(t, s, idx)
	assert.Nil(t, it)

	it = sf.Seek(func(i int) bool {
		return i == 4
	})
	i, done := it.Cur()
	assert.Equal(t, 4, i)
	assert.False(t, done)

	idx = it.Idx()
	it = iter.Seek(it, fn)
	assert.Len(t, s, idx)
	assert.Nil(t, it)
}

func TestFactoryFor(t *testing.T) {
	var sf iter.Factory[byte] = sliceFactory([]byte("hello"))
	out := ""
	fn := func(b byte) {
		out += string(b)
	}
	sf.For(fn)
	assert.Equal(t, "hello", out)
}

func TestFactoryForIdx(t *testing.T) {
	s := []int{3, 1, 4, 1, 5, 9}
	var sf iter.Factory[int] = sliceFactory(s)

	fn := func(i, idx int) {
		assert.Equal(t, s[idx], i)
	}
	c := sf.ForIdx(fn)
	assert.Len(t, s, c)
}

func TestFactoryConcurrent(t *testing.T) {
	s := []int{3, 1, 4, 1, 5, 9}
	var sf iter.Factory[int] = sliceFactory(s)
	var c int32
	wg := sf.Concurrent(func(i, idx int) {
		assert.Equal(t, s[idx], i)
		atomic.AddInt32(&c, 1)
	})

	wg.Wait()
	assert.Len(t, s, int(c))
}

func TestFactoryChannel(t *testing.T) {
	s := []int{3, 1, 4, 1, 5, 9}
	var sf iter.Factory[int] = sliceFactory(s)

	idx := 0
	for i := range sf.Channel(0) {
		assert.Equal(t, s[idx], i)
		idx++
	}
	assert.Len(t, s, idx)
}

func TestFactoryWrap(t *testing.T) {
	s := []int{3, 1, 4, 1, 5, 9}
	var sf iter.Factory[int] = sliceFactory(s)
	w, i, done := sf.Wrap()
	assert.Equal(t, s[0], i)
	assert.False(t, done)
	w.ForIdx(func(i, idx int) {
		assert.Equal(t, s[idx], i)
	})
}

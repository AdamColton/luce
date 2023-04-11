package liter_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/util/liter"
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
	it := liter.Seek[int](s, fn)
	assert.Len(t, s.Slice, idx)
	assert.Nil(t, it)

	s.idx = 0
	it = liter.Seek[int](s, func(i int) bool {
		return i == 4
	})
	i, done := it.Cur()
	assert.Equal(t, 4, i)
	assert.False(t, done)
	idx = s.idx

	it = liter.Seek[int](s, fn)
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
	it := liter.Seek[int](s, fn)

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
	liter.For[byte](si, fn)
	assert.Equal(t, expected, out)
}

func ExampleFor() {
	si := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	fn := func(i int) {
		fmt.Println(i)
	}
	liter.For[int](si, fn)
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
	c := liter.ForIdx[int](si, fn)
	assert.Len(t, s, c)
}

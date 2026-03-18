package morph

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
)

// Slice applies the Val to every value in a and returns the values that are
// included as a slice. The buffer is used if it has sufficient capacity.
func (vt Val[In, Out]) Slice(a []In, buf []Out) slice.Slice[Out] {
	i := vt.SliceIter(a)
	return slice.FromIter(i, buf)
}

// SliceIter returns an iterator over the included values of a after the Val is
// applied.
func (vt Val[In, Out]) SliceIter(a []In) liter.Wrapper[Out] {
	return vt.Iter(slice.New(a).Iter())
}

// IterSlice applies the Val to the values of the iterator, from its current
// value, and returns the values that are included as a slice. The buffer is used
// if it has sufficient capacity.
func (vt Val[In, Out]) IterSlice(a liter.Iter[In], buf []Out) slice.Slice[Out] {
	return slice.FromIter(vt.Iter(a), buf)
}

// Slice applies the ValAll to every value in a and returns the result as a
// slice. The buffer is used if it has sufficient capacity.
func (vt ValAll[In, Out]) Slice(a []In, buf []Out) slice.Slice[Out] {
	i := vt.SliceIter(a)
	return slice.FromIter(i, buf)
}

// SliceIter returns an iterator over the values of a after the ValAll is
// applied.
func (vt ValAll[In, Out]) SliceIter(a []In) liter.Wrapper[Out] {
	return vt.Iter(slice.New(a).Iter())
}

// IterSlice applies the ValAll to the values of the iterator, from its current
// value, and returns the result as a slice. The buffer is used if it has
// sufficient capacity.
func (vt ValAll[In, Out]) IterSlice(a liter.Iter[In], buf []Out) slice.Slice[Out] {
	return slice.FromIter(vt.Iter(a), buf)
}

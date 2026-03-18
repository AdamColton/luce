package morph

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
)

func (vt Val[In, Out]) Slice(a []In, buf []Out) slice.Slice[Out] {
	i := vt.SliceIter(a)
	return slice.FromIter(i, buf)
}

func (vt Val[In, Out]) SliceIter(a []In) liter.Wrapper[Out] {
	return vt.Iter(slice.New(a).Iter())
}

func (vt Val[In, Out]) IterSlice(a liter.Iter[In], buf []Out) slice.Slice[Out] {
	return slice.FromIter(vt.Iter(a), buf)
}

func (vt ValAll[In, Out]) Slice(a []In, buf []Out) slice.Slice[Out] {
	i := vt.SliceIter(a)
	return slice.FromIter(i, buf)
}

func (vt ValAll[In, Out]) SliceIter(a []In) liter.Wrapper[Out] {
	return vt.Iter(slice.New(a).Iter())
}

func (vt ValAll[In, Out]) IterSlice(a liter.Iter[In], buf []Out) slice.Slice[Out] {
	return slice.FromIter(vt.Iter(a), buf)
}

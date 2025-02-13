package slice

import "github.com/adamcolton/luce/util/liter"

// TransformFunc converts one type to another and can be applied to an iterator.
// TransformFunc is also the same as liter.TransformFunc.
type TransformFunc[In, Out any] func(In, int) (out Out, include bool)

// Factory creates a slice by transforming the values from the iterator returned
// by the Factory. This uses a recursive call, so if it might blow the stack, do
// it in a for loop.
func (fn TransformFunc[In, Out]) Factory(f liter.Factory[In], buf []Out) (out Slice[Out]) {
	return FromIterFactory(liter.TransformFunc[In, Out](fn).Factory(f), buf)
}

// Transform creates a slice by transforming the values from the iterator.
// This uses a recursive call, so if it might blow the stack, do it in a
// for loop.
func (fn TransformFunc[In, Out]) Transform(in liter.Iter[In], buf []Out) Slice[Out] {
	return FromIter(liter.TransformFunc[In, Out](fn).New(in), buf)
}

// TransformSlice uses the TransformFunc to convert the values from one slice
// to another for all the included values.
func (fn TransformFunc[In, Out]) Slice(in []In, buf []Out) Slice[Out] {
	return fn.Transform(NewIter(in), buf)
}

// Transform applies the TransformFunc to an iterator and collects all the
// included results in a slice.
func Transform[In, Out any](in liter.Iter[In], buf []Out, fn TransformFunc[In, Out]) Slice[Out] {
	return fn.Transform(in, buf)
}

// Transform one slice to another. The transformation function's second return
// is a bool indicating if the returned value should be included in the result.
// The returned Slice is sized exactly to the output.
func TransformSlice[In, Out any](in []In, buf []Out, fn TransformFunc[In, Out]) Slice[Out] {
	return fn.Transform(NewIter(in), buf)
}

// ForAll is a helper function for transformers. Transformers have an index
// argument and a bool return that are often not used.
func ForAll[In, Out any](fn func(in In) Out) TransformFunc[In, Out] {
	return func(in In, idx int) (Out, bool) {
		return fn(in), true
	}
}

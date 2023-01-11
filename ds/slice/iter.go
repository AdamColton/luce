package slice

import (
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/upgrade"
)

// == projects.Code.luce.slice ==
// [ ] order slice by slice
//	given 2 slices, apply SortInPlace
//	to the first, based on if a value appears in the second
//	but it should be able to swap if included goes to the front or back
//	and it should iterate over the shorter and compare to the longer
//	because the runtime will be m*log(n)
//	see Delaunay.missingEdges for inspiration

// Iter wraps a slice to fulfill liter.Iter.
type Iter[T any] struct {
	// S is the slice being iterated over.
	S []T
	// I is the current index.
	I int
	// DoneFn determines when the iterator is done. If it does not check the
	// bounds, an out of bounds panic is likely.
	DoneFn func(*Iter[T]) bool
}

// NewIter creates an Iter for iterating over the slice. Done defaults to a func
// that compares to the length of the slice.
func NewIter[T any](s []T) liter.Wrapper[T] {
	return liter.Wrap(&Iter[T]{
		S: s,
		DoneFn: func(i *Iter[T]) bool {
			return i.I >= len(i.S)
		},
	})
}

// Cur fulfills liter.Iter. It returns the value at the current index and the
// done bool. If it is done, it will return the zero value for the type.
func (i *Iter[T]) Cur() (t T, done bool) {
	done = i.Done()
	if !done {
		t = i.S[i.I]
	}
	return
}

// Start fulfills liter.Starter. It sets the index to zero and returns the first
// element in the slice. If it is done, it will return the zero value for the
// type.
func (i *Iter[T]) Start() (t T, done bool) {
	i.I = 0
	return i.Cur()
}

// Next fulfills liter.Iter. It increments the index and returns the value at
// that index and the done bool. If it is done, it will return the zero value
// for the type.
func (i *Iter[T]) Next() (t T, done bool) {
	i.I++
	return i.Cur()
}

// Done fulfills liter.Iter. It calls the underlying DoneFn.
func (i *Iter[T]) Done() bool {
	return i.DoneFn(i)
}

// Idx fulfills liter.Iter. It returns the current value of I.
func (i *Iter[T]) Idx() int {
	return i.I
}

// Slice fulfills Slicer. If a buffer with capacity is provided, the whole
// slice is copied into it (regardless of the current index). If not, the
// underlying slice is returned without copying, so it shares memory with the
// iterator.
func (i *Iter[T]) Slice(buf []T) []T {
	if cap(buf) > 0 {
		return append(buf[:0], i.S...)
	}
	return i.S
}

// IterFactory fulfills liter.Factory.
func IterFactory[T any](s []T) liter.Factory[T] {
	return func() (i liter.Iter[T], t T, done bool) {
		i = NewIter(s)
		t, done = i.Cur()
		return
	}
}

// FromIter creates a slice from a liter.Iter. If the Iter fulfills Slicer
// (possibly through a Wrapper), then the Slice method will be used. If it
// fulfills Lener directly, that will be used for buffer allocation. Otherwise the values
// are collected recursively, so the stack grows with the length of the Iter.
// An Iter that is done returns nil.
func FromIter[T any](i liter.Iter[T], buf []T) Slice[T] {
	if s, ok := upgrade.To[Slicer[T]](i); ok {
		return s.Slice(buf)
	}

	if ln, ok := i.(Lener); ok {
		agg := Buffer[T](buf).Empty(ln.Len())
		return liter.Appender[T]().Iter(agg, i)
	}

	t, done := i.Cur()
	if done {
		return nil
	}
	size := 1
	out := recurFromIter(i, buf, &size)
	out[0] = t
	return out
}

func recurFromIter[T any](i liter.Iter[T], buf []T, size *int) (out Slice[T]) {
	t, done := i.Next()
	if done {
		out = Buffer[T](buf).Slice(*size)
	} else {
		*size++
		out = recurFromIter(i, buf, size)
		*size--
		out[*size] = t
	}
	return
}

// FromIterFactory creates a slice from a liter.Factory, as FromIter does. If
// the Factory is done, nil is returned.
func FromIterFactory[T any](f liter.Factory[T], buf []T) Slice[T] {
	i, _, done := f()
	if done {
		return nil
	}
	return FromIter(i, buf)
}

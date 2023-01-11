package slice

import (
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/upgrade"
)

// Iter wraps a slice to fulfillliter.Iter.
type Iter[T any] struct {
	S []T
	I int
	// DoneFn determines when the iterator is done. If it does not check the
	// bounds, an out of bounds panic is likely.
	DoneFn func(*Iter[T]) bool
}

// NewIter creates an Iter for iterating over the slice. Done defaults to a func
// that compares to the length of the slice.
func NewIter[T any](s []T) liter.Wrapper[T] {
	return liter.Wrapper[T]{&Iter[T]{
		S: s,
		DoneFn: func(i *Iter[T]) bool {
			return i.I >= len(i.S)
		},
	}}
}

// Cur fulfillsliter.Iter. It returns the value at the current index and the done bool. If it is done,
// it will return the zero value for the type.
func (i *Iter[T]) Cur() (t T, done bool) {
	done = i.Done()
	if !done {
		t = i.S[i.I]
	}
	return
}

// Start fulfillsliter.Iter. It sets the index to zero and returns the first element in the slice. If
// it is done, it will return the zero value for the type.
func (i *Iter[T]) Start() (t T, done bool) {
	i.I = 0
	return i.Cur()
}

// Next fulfillsliter.Iter. It increments the index and returns the value at that index and the done
// bool. If it is done, it will return the zero value for the type.
func (i *Iter[T]) Next() (t T, done bool) {
	i.I++
	return i.Cur()
}

// Done fulfillsliter.Iter. It calls the underlying DoneFn.
func (i *Iter[T]) Done() bool {
	return i.DoneFn(i)
}

// Idx fulfillsliter.Iter. It returns the current value of I.
func (i *Iter[T]) Idx() int {
	return i.I
}

// Slice fulfills Slicer. If a buffer is provided, a copy is made. If no buffer
// is provided, the underlying slice is returned.
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

// IterSlice creates a slice from an liter.Iter. If the Iter fulfills Slicer,
// then the Slice method will be used. If it fulfill Lener, that will be used
// for buffer allocation.
func IterSlice[T any](i liter.Iter[T], buf []T) []T {
	if s, ok := upgrade.To[Slicer[T]](i); ok {
		return s.Slice(buf)
	}

	agg := Buffer[T](buf).Lener(i)
	return liter.Appender[T]().Iter(agg, i)
}

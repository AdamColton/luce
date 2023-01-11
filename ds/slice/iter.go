package slice

import (
	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/upgrade"
)

// Iter wraps a slice to fulfill iter.Iter.
type Iter[T any] struct {
	S []T
	I int
	// DoneFn determines when the iterator is done. If it does not check the
	// bounds, an out of bounds panic is likely.
	DoneFn func(*Iter[T]) bool
}

// NewIter creates an Iter for iterating over the slice. Done defaults to a func
// that compares to the length of the slice.
func NewIter[T any](s []T) *Iter[T] {
	return &Iter[T]{
		S: s,
		DoneFn: func(i *Iter[T]) bool {
			return i.I >= len(i.S)
		},
	}
}

// Cur fulfills iter.Iter. It returns the value at the current index and the done bool. If it is done,
// it will return the zero value for the type.
func (i *Iter[T]) Cur() (t T, done bool) {
	done = i.Done()
	if !done {
		t = i.S[i.I]
	}
	return
}

// Start fulfills iter.Iter. It sets the index to zero and returns the first element in the slice. If
// it is done, it will return the zero value for the type.
func (i *Iter[T]) Start() (t T, done bool) {
	i.I = 0
	return i.Cur()
}

// Next fulfills iter.Iter. It increments the index and returns the value at that index and the done
// bool. If it is done, it will return the zero value for the type.
func (i *Iter[T]) Next() (t T, done bool) {
	i.I++
	return i.Cur()
}

// Done fulfills iter.Iter. It calls the underlying DoneFn.
func (i *Iter[T]) Done() bool {
	return i.DoneFn(i)
}

// Idx fulfills iter.Iter. It returns the current value of I.
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

// IterFactory fulfills iter.Factory.
func IterFactory[T any](s []T) iter.Factory[T] {
	return func() (i iter.Iter[T], t T, done bool) {
		i = NewIter(s)
		t, done = i.Cur()
		return
	}
}

// IterSlice creates a slice from an iter.Iter. If the Iter fulfills Slicer,
// then the Slice method will be used. If it fulfill Lener, that will be used
// for buffer allocation.
func IterSlice[T any](i iter.Iter[T], buf []T) []T {
	var s Slicer[T]
	if upgrade.Upgrade(i, &s) {
		return s.Slice(buf)
	}

	out := BufferLener(i, buf)
	iter.For(i, func(t T, idx int) {
		out = append(out, t)
	})

	return out
}

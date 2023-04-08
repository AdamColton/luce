package list

import (
	"github.com/adamcolton/luce/ds/slice"
)

// Reverse a List
type Reverse[T any] struct {
	List[T]
}

// NewReverse creates a Wrapped list that is the reverse of l.
func NewReverse[T any](l List[T]) Wrapper[T] {
	return Wrapper[T]{&Reverse[T]{l}}
}

// Note that Reverse should not have an Upgrade because it is modifying the
// underlying List. For instance, if it was to be upgraded to Stringer and the
// underlying List fulfilled Stringer, it would do so without reversing the
// list.

// AtIdx returns the value at the index. Fulfills List.
func (r Reverse[T]) AtIdx(idx int) T {
	return r.List.AtIdx(r.Len() - 1 - idx)
}

// Slice fulfills Slicer. It converts the underlying Reversed List to a slice.
func (r Reverse[T]) Slice(buf []T) []T {
	ln := r.Len()
	out := slice.Buffer[T](buf).Empty(ln)
	ln--
	for i := 0; i <= ln; i++ {
		out = append(out, r.List.AtIdx(ln-i))
	}
	return out
}

// Wrap the SliceList to add Wrapper methods.
func (r Reverse[T]) Wrap() Wrapper[T] {
	return Wrapper[T]{r}
}

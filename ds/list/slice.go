package list

import "github.com/adamcolton/luce/ds/slice"

// SliceList wraps a slice to fulfill List.
type SliceList[T any] []T

// AtIdx returns the value at the index.
func (sl SliceList[T]) AtIdx(idx int) T {
	return sl[idx]
}

// Len returns the length of the underlying slice.
func (sl SliceList[T]) Len() int {
	return len(sl)
}

// Slice fulfills Slicer and returns the underlying slice.
func (sl SliceList[T]) Slice(buf []T) []T {
	return sl
}

// ToSlice converts a List to slice. If the underlying List implements Slicer,
// that will be invoked.
func ToSlice[T any](l List[T], buf []T) []T {
	return slice.IterSlice[T](NewIter(l), buf)
}

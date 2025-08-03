package list

import (
	"github.com/adamcolton/luce/math/ints"
)

// ByIdx uses an index list to remap the indicies of the Source list. ByIdx
// fulfills List.
type ByIdx[T any, N ints.Number] struct {
	Source List[T]
	Idxs   List[N]
}

// NewByIdx creates a new ByIdx
func NewByIdx[T any, N ints.Number](src List[T], idxs List[N]) ByIdx[T, N] {
	return ByIdx[T, N]{
		Source: src,
		Idxs:   idxs,
	}
}

func (i ByIdx[T, N]) Wrap() Wrapper[T] {
	return Wrapper[T]{i}
}

// Len returns the length of the index list.
func (i ByIdx[T, N]) Len() int {
	return i.Idxs.Len()
}

// AtIdx loops up idx in i.Idxs to find the correspeding index in i.Source.
func (i ByIdx[T, N]) AtIdx(idx int) T {
	ii := int(i.Idxs.AtIdx(idx))
	return i.Source.AtIdx(ii)
}

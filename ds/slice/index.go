package slice

import (
	"github.com/adamcolton/luce/math/ints"
)

// Index represents a range within a slice. Negative numbers are interpreted
// relative to the end of a slice. Index fulfills list.List.
type Index [2]int

// Sub uses the Index to return a sub slice.
func (s Slice[T]) Sub(i Index) Slice[T] {
	i, _ = i.Size(len(s))
	return s[i[0]:i[1]]
}

// NewIndex creates an Index from a start and a length. To create an Index from
// a start and an end, instanciate it directly.
func NewIndex(start, ln int) Index {
	return Index{start, start + ln}
}

// Next creates an Index that starts immidatly after index i.
func (i Index) Next(ln int) Index {
	return Index{i[1], i[1] + ln}
}

// IdxMake creates a Slice sized to the second value in the Index.
func IdxMake[T any](idx Index) Slice[T] {
	return make([]T, idx[1])
}

// Last returns the last value the Index includes, which is one less than i[1].
func (i Index) Last() int {
	return i[1] - 1
}

// Len returns the length of the Index.
func (i Index) Len() int {
	return i[1] - i[0]
}

// AtIdx fulfills list.List. It returns the index relative to
func (i Index) AtIdx(idx int) int {
	return i[0] + idx
}

// Size the index to a known length if it contains values relative to the end.
func (i Index) Size(ln int) (Index, bool) {
	a, aOk := ints.Idx(i[0], ln)
	b, bOk := ints.Idx(i[1], ln)
	return Index{a, b}, aOk && bOk
}

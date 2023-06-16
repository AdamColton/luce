// Package huffslice compresses a slice using Huffman encoding.
package huffslice

import (
	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/util/liter"
)

// Slice is a slice compressed with Huffman encoding. Values that occur only once
// (and SingleToken itself) are not worth a place in the Tree, so they are stored
// in Singles and replaced by SingleToken in the encoded bits. Create one with an
// Encoder and read it back with Iter.
type Slice[T comparable] struct {
	huffman.Tree[T]
	// Encoded holds the Huffman encoded values.
	Encoded *rye.Bits
	// Singles holds the values that occur a single time, in order. They are
	// replaced with SingleToken in Encoded and stored separately from the Tree.
	Singles slice.Slice[T]
	// SingleToken stands in for the singles in Encoded. Any value that equals it
	// is treated as a single.
	SingleToken T
}

// Encoder is used to create the Slice. Fill Slice with the values, then call
// Encode.
type Encoder[T comparable] struct {
	// Slice holds the values to be encoded.
	Slice slice.Slice[T]
	// SingleToken is the value that replaces the values that occur once.
	SingleToken T
}

// NewEncoder creates an Encoder with an empty Slice with a capacity of ln and
// sets the value for the single token.
func NewEncoder[T comparable](ln int, singleToken T) *Encoder[T] {
	return &Encoder[T]{
		Slice:       make([]T, 0, ln),
		SingleToken: singleToken,
	}
}

// Encode the values in the Encoder's Slice to a Slice. The Tree is built from the
// values that occur more than once, plus SingleToken. If every value is a
// single, the Tree holds only SingleToken, its code has no bits, and Encoded is
// empty; Iter then reads the Singles directly.
func (e *Encoder[T]) Encode() *Slice[T] {
	freq := make(map[T]int, len(e.Slice))
	for _, t := range e.Slice {
		freq[t]++
	}

	singlesCount := 0
	for t, c := range freq {
		if c == 1 || t == e.SingleToken {
			singlesCount++
			delete(freq, t)
		}
	}
	freq[e.SingleToken] = singlesCount

	out := &Slice[T]{
		Tree:        huffman.MapNew(freq),
		SingleToken: e.SingleToken,
		Singles:     slice.Make[T](0, singlesCount),
	}

	fnList := list.Generator[T]{
		Length: len(e.Slice),
		Fn: func(idx int) T {
			t := e.Slice[idx]
			if freq[t] > 1 && t != e.SingleToken {
				return t
			}
			out.Singles = append(out.Singles, t)
			return e.SingleToken
		},
	}

	out.Encoded = huffman.Encode[T](fnList, huffman.NewLookup(out.Tree))

	return out
}

// Iter creates an iterator for decoding the Slice. It returns the values in
// their original order, starting from the first. It shares Encoded, so it moves
// Encoded.Idx.
func (s *Slice[T]) Iter() liter.Iter[T] {
	if s.Encoded.Ln == 0 {
		return slice.NewIter(s.Singles)
	}
	i := &sliceiter[T]{
		Iter:        s.Tree.Iter(s.Encoded),
		singleToken: s.SingleToken,
		singles:     s.Singles,
	}
	i.Start()
	return i
}

type sliceiter[T comparable] struct {
	liter.Iter[T]
	sIdx        int
	singleToken T
	singles     slice.Slice[T]
}

func (si *sliceiter[T]) Start() (t T, done bool) {
	si.sIdx = -1
	t, done = si.Iter.(liter.Starter[T]).Start()
	if !done && t == si.singleToken {
		si.sIdx++
		t = si.singles[si.sIdx]
	}
	return
}

func (si *sliceiter[T]) Next() (t T, done bool) {
	t, done = si.Iter.Next()
	if !done && t == si.singleToken {
		si.sIdx++
		t = si.singles[si.sIdx]
	}
	return
}

func (si *sliceiter[T]) Cur() (t T, done bool) {
	t, done = si.Iter.Cur()
	if !done && t == si.singleToken {
		t = si.singles[si.sIdx]
	}
	return
}

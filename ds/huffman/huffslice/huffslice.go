// Package huffslice compresses a slice using Huffman encoding.
package huffslice

import (
	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/util/iter"
)

// Slice encoded using Huffman compression
type Slice[T comparable] struct {
	huffman.Tree[T]
	Encoded *rye.Bits
	// Values that occure a single time are replaced wit SingleToken and stored
	// speratly from the Tree.
	Singles     slice.Slice[T]
	SingleToken T
}

// Encoder is used to create the Slice.
type Encoder[T comparable] struct {
	Slice       slice.Slice[T]
	SingleToken T
}

// NewEncoder creates a slice with a capacity of ln and sets the values for the
// single token.
func NewEncoder[T comparable](ln int, singleToken T) *Encoder[T] {
	return &Encoder[T]{
		Slice:       make([]T, 0, ln),
		SingleToken: singleToken,
	}
}

// Encode to Slice.
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

// Iter creates an iterator for decoding the Slice.
func (s *Slice[T]) Iter() iter.Iter[T] {
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
	iter.Iter[T]
	sIdx        int
	singleToken T
	singles     slice.Slice[T]
}

func (si *sliceiter[T]) Start() (t T, done bool) {
	si.sIdx = -1
	t, done = si.Iter.(iter.Starter[T]).Start()
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

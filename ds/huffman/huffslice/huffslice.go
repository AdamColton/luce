// Package huffslice compresses a slice using Huffman encoding.
package huffslice

import (
	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye"
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

// Decode from a huffslice.Slice.
func (s *Slice[T]) Decode() slice.Slice[T] {
	sIdx := 0
	out := slice.New(s.Tree.Iter(s.Encoded).Slice(nil))
	if len(out) == 0 {
		return s.Singles
	}
	out.Iter().ForIdx(func(t T, idx int) {
		if t == s.SingleToken {
			out[idx] = s.Singles[sIdx]
			sIdx++
		}
	})
	return out
}

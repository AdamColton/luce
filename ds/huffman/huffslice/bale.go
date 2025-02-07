package huffslice

import (
	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/serial/rye"
)

// SliceBale is a flat form of a Slice that can be serialized.
type SliceBale[T comparable] struct {
	// TreeBale is the bale of the Slice's Tree.
	*huffman.TreeBale[T]
	// Encoded, Singles and SingleToken are the fields of the Slice.
	Encoded     *rye.Bits
	Singles     []T
	SingleToken T
}

// Bale returns the Slice as a SliceBale.
func (s *Slice[T]) Bale() *SliceBale[T] {
	return &SliceBale[T]{
		TreeBale:    s.Tree.Bale(),
		Encoded:     s.Encoded,
		Singles:     s.Singles,
		SingleToken: s.SingleToken,
	}
}

// UnbaleTo replaces the contents of s with the Slice the bale describes.
func (bale *SliceBale[T]) UnbaleTo(s *Slice[T]) {
	s.Tree = bale.TreeBale.Unbale()
	s.Encoded = bale.Encoded
	s.Singles = bale.Singles
	s.SingleToken = bale.SingleToken
}

// Unbale creates a Slice from the bale.
func (bale *SliceBale[T]) Unbale() *Slice[T] {
	out := &Slice[T]{}
	bale.UnbaleTo(out)
	return out
}

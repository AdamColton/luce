package huffslice

import (
	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/serial/rye"
)

type SliceBale[T comparable] struct {
	*huffman.TreeBale[T]
	Encoded     *rye.Bits
	Singles     []T
	SingleToken T
}

func (s *Slice[T]) Bale() *SliceBale[T] {
	return &SliceBale[T]{
		TreeBale:    s.Tree.Bale(),
		Encoded:     s.Encoded,
		Singles:     s.Singles,
		SingleToken: s.SingleToken,
	}
}

func (bale *SliceBale[T]) UnbaleTo(s *Slice[T]) {
	s.Tree = bale.TreeBale.Unbale()
	s.Encoded = bale.Encoded
	s.Singles = bale.Singles
	s.SingleToken = bale.SingleToken
}

func (bale *SliceBale[T]) Unbale() *Slice[T] {
	out := &Slice[T]{}
	bale.UnbaleTo(out)
	return out
}

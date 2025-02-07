package huffman

import "github.com/adamcolton/luce/math/ints"

type TreeBale[T any] struct {
	Values []T
	Nodes  [][2]uint32
}

func (t *Tree[T]) Bale() *TreeBale[T] {
	tb := &TreeBale[T]{
		Values: make([]T, 0, t.ln),
		Nodes:  make([][2]uint32, 0, 2*t.ln),
	}
	tb.addNode(t.huffNode)
	return tb
}

func (bale *TreeBale[T]) UnbaleTo(t *Tree[T]) {
	t.ln = len(bale.Values)
	t.huffNode = bale.getNode(0)
}

func (bale *TreeBale[T]) Unbale() *Tree[T] {
	out := &Tree[T]{}
	bale.UnbaleTo(out)
	return out
}

func (bale *TreeBale[T]) getNode(idx uint32) *huffNode[T] {
	n := bale.Nodes[idx]
	hn := &huffNode[T]{}
	if n[0] == ints.MaxU32 {
		hn.v = bale.Values[n[1]]
	} else {
		hn.branch[0] = bale.getNode(n[0])
		hn.branch[1] = bale.getNode(n[1])
	}
	return hn
}

func (bale *TreeBale[T]) addNode(hn *huffNode[T]) uint32 {
	idx := uint32(len(bale.Nodes))
	var n [2]uint32
	bale.Nodes = append(bale.Nodes, n)
	if hn.branch[0] == nil {
		n[0] = ints.MaxU32
		n[1] = uint32(len(bale.Values))
		bale.Values = append(bale.Values, hn.v)
	} else {
		n[0] = bale.addNode(hn.branch[0])
		n[1] = bale.addNode(hn.branch[1])
	}
	bale.Nodes[idx] = n
	return idx
}

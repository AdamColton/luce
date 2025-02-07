package huffman

import "github.com/adamcolton/luce/math/ints"

// TreeBale is a flat form of a Tree that can be serialized. Node 0 is the root.
// A branch is a pair of indexes into Nodes. A leaf is a pair of math.MaxUint32
// and an index into Values.
type TreeBale[T any] struct {
	// Values are the values of the leaves.
	Values []T
	// Nodes are the nodes of the tree.
	Nodes [][2]uint32
}

// Bale returns the Tree as a TreeBale.
func (t *Tree[T]) Bale() *TreeBale[T] {
	tb := &TreeBale[T]{
		Values: make([]T, 0, t.ln),
		Nodes:  make([][2]uint32, 0, 2*t.ln),
	}
	tb.addNode(t.huffNode)
	return tb
}

// UnbaleTo replaces the contents of t with the Tree the bale describes.
func (bale *TreeBale[T]) UnbaleTo(t *Tree[T]) {
	t.ln = len(bale.Values)
	t.huffNode = bale.getNode(0)
}

// Unbale creates a Tree from the bale.
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

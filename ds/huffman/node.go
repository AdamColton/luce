package huffman

import "github.com/adamcolton/luce/serial/rye"

type huffNode[T any] struct {
	branch [2]*huffNode[T]
	v      T
}

func (d *Frequency[T]) root() *root[T] {
	return newLeaf(d.Val, d.Count)
}

type root[T any] struct {
	node *huffNode[T]
	sum  int
}

func newLeaf[T any](v T, sum int) *root[T] {
	return &root[T]{
		node: &huffNode[T]{
			v: v,
		},
		sum: sum,
	}
}

func newBranch[T any](n0, n1 *huffNode[T], sum int) *root[T] {
	return &root[T]{
		node: &huffNode[T]{
			branch: [2]*huffNode[T]{n0, n1},
		},
		sum: sum,
	}
}

// Read from b until a leaf is encountered, return the leaf value.
func (n *huffNode[T]) Read(b *rye.Bits) T {
	if n.branch[0] == nil {
		return n.v
	}
	return n.branch[b.Read()].Read(b)
}

// ReadAll bits, traversing the Huffman tree.
func (n *huffNode[T]) ReadAll(b *rye.Bits) []T {
	var out []T
	for b.Idx < b.Ln {
		out = append(out, n.Read(b))
	}
	return out
}

func (n *huffNode[T]) All(fn func(T)) {
	if n.branch[0] == nil {
		fn(n.v)
	} else {
		n.branch[0].All(fn)
		n.branch[1].All(fn)
	}
}

package merkle

import "hash"

type branch struct {
	children     [2]node
	digest, data []byte
	leaves       int
}

func (b *branch) update(h hash.Hash) {
	h.Reset()
	h.Write(b.children[0].Digest())
	h.Write(b.children[1].Digest())
	b.digest = h.Sum(b.digest[:0])
	b.leaves = b.children[0].Leaves() + b.children[1].Leaves()
}

func (b *branch) Digest() []byte {
	return b.digest
}

func (b *branch) Data() []byte {
	return b.data
}

func (b *branch) Leaves() int {
	return b.leaves
}

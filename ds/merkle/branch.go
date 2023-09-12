package merkle

import "hash"

type branch struct {
	children     [2]node
	digest, data []byte
	leaves       int
}

func (b *branch) update(h hash.Hash) (ln, leaves int) {
	ln0, leaves0 := b.children[0].update(h)
	ln1, leaves1 := b.children[1].update(h)
	h.Reset()
	h.Write(b.children[0].Digest())
	h.Write(b.children[1].Digest())
	b.digest = h.Sum(b.digest[:0])
	b.leaves = leaves0 + leaves1
	return ln0 + ln1, b.leaves
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

func (b *branch) Len() int {
	return len(b.data)
}

func (b *branch) stitchData(in []byte) int {
	ln := b.children[0].stitchData(in)
	ln += b.children[1].stitchData(in[ln:])
	b.data = in[:ln]
	return ln
}

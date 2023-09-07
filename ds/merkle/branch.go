package merkle

import "hash"

type branch struct {
	children     [2]node
	digest, data []byte
}

func (b *branch) setDigest(h hash.Hash) {
	h.Reset()
	h.Write(b.children[0].Digest())
	h.Write(b.children[1].Digest())
	b.digest = h.Sum(b.digest[:0])
}

func (b *branch) Digest() []byte {
	return b.digest
}

func (b *branch) Data() []byte {
	return b.data
}

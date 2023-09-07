package merkle

import (
	"hash"
	"math/bits"
)

type builder struct {
	maxLeafSize int
	h           func() hash.Hash
}

// NewBuilder creates a Builder whose trees will follow the limits set by
// maxSize and branch and will use the provided hash.
func NewBuilder(maxSize uint32, h func() hash.Hash) Builder {
	return builder{
		maxLeafSize: int(maxSize),
		h:           h,
	}
}

func (b builder) Build(data []byte) Tree {
	if len(data) == 0 {
		return nil
	}

	t := &tree{
		h: b.h(),
	}
	t.node, t.leaves = makeTree(0, b.maxLeafSize, data, t.h)
	t.depth = uint32(bits.Len(uint(t.leaves)))
	return t
}

func makeTree(idx uint32, maxLeafSize int, data []byte, h hash.Hash) (node, uint32) {
	ln := len(data)
	if ln <= maxLeafSize {
		return newDataLeaf(data, idx, h), idx + 1
	}

	leaves := divUp(ln, maxLeafSize)
	split := (leaves / 2) * maxLeafSize

	n := &branch{
		data: data,
	}
	n.children[0], idx = makeTree(idx, maxLeafSize, data[:split], h)
	n.children[1], idx = makeTree(idx, maxLeafSize, data[split:], h)
	n.update(h)
	return n, idx
}

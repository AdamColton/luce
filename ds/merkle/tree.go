package merkle

import "hash"

type tree struct {
	node
	h      hash.Hash
	leaves uint32
}

func (t *tree) Leaves() int {
	return int(t.leaves)
}

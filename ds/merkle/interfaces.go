package merkle

import "hash"

// Builder handles the logic of generating and populating a tree from data.
type Builder interface {
	Build(data []byte) Tree
}

// node that has been populated with data and hashes
type node interface {
	// Digest returns the digest of the Tree for validation
	Digest() []byte
	// Data returns the entire data of the tree
	Data() []byte
	// Leaves is the total count of leaves.
	Leaves() int
	// Len of the underlying []byte created by joining all the leaves.
	Len() int
	update(hash.Hash) (ln, leaves int)
	stitchData([]byte) int
}

type Tree interface {
	node
	Leaf(int) *Leaf
	Description() Description
}

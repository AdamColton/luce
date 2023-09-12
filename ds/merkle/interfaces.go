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

// Tree is a Merkle tree over a []byte of data. It is also an io.Reader and an
// io.Seeker over that data.
type Tree interface {
	node
	// Leaf returns the Leaf at idx, with the rows that validate it, or nil if idx is
	// out of range.
	Leaf(int) *Leaf
	// Description returns what an Assembler needs to validate the Leaves of this
	// tree: its digest and the number of leaves.
	Description() Description
}

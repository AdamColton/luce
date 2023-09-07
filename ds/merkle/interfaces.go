package merkle

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
	Leaves() int
}

// Tree is a Merkle tree over a []byte of data. It is also an io.Reader and an
// io.Seeker over that data.
type Tree interface {
	node
	// Leaf returns the Leaf at idx, with the rows that validate it, or nil if idx is
	// out of range.
	Leaf(int) *Leaf
}

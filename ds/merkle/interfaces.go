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

type Tree interface {
	node
	Leaf(int) *Leaf
}

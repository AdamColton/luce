package merkle

type node interface {
	Digest() []byte
	Count() uint32
	Depth() int
	size() int
	maxIdx() uint32
	have([]uint32) []uint32
}

// Builder handles the logic of generating and populating a tree from data.
type Builder interface {
	Build(data []byte) Tree
}

// Tree that has been populated with data and hashes
type Tree interface {
	// Digest returns the digest of the Tree for validation
	Digest() []byte
	// Data returns the entire data of the tree
	Data() []byte
	// Leaf represents the data of a single leaf and includes all Uncle
	// hashes allowing the data to be validated against the tree hash.
	Leaf(idx uint32) Leaf
	// Maximum Depth of the tree
	Depth() int
	// Count of the leaf nodes
	Count() uint32
	size() int

	// Creates a Description of the Tree
	Description() Description
	// Read fulfills io.Reader
	Read(p []byte) (n int, err error)
}

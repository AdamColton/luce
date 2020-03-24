package merkle

import "hash"

// Description of tree.
type Description struct {
	Digest []byte
	Count  uint32
}

func (b *branch) Description() Description {
	return Description{
		Digest: b.Digest(),
		Count:  uint32(b.Count()),
	}
}

// Assembler creates a new Assembler from the Description. All Leaves added to
// the tree will be validated against against the initial validator.
func (d Description) Assembler(h hash.Hash) *Assembler {
	return &Assembler{
		d: d,
		h: h,
	}
}

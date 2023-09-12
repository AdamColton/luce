package merkle

import (
	"hash"

	"github.com/adamcolton/luce/serial/rye"
)

// Description of tree.
type Description struct {
	Digest []byte
	Leaves uint32
}

// Assembler creates a new Assembler from the Description. All Leaves added to
// the tree will be validated against against the initial validator.
func (d Description) Assembler(h hash.Hash) *Assembler {
	bits := d.Leaves / 8
	if d.Leaves%8 != 0 {
		bits++
	}
	return &Assembler{
		d:         d,
		h:         h,
		remaining: int(d.Leaves),
		bools: &rye.Bits{
			Data: make([]byte, bits),
			Ln:   int(d.Leaves),
		},
	}
}

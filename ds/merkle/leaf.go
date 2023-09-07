package merkle

import (
	"hash"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/serial/rye"
)

// ValidatorRow represents a Branch in a Merkle Tree. It provides the digests
// needed for validation.
type ValidatorRow struct {
	SiblingDigest []byte
	IsFirst       bool
}

// Leaf represents a Data Leaf in a Merkle Tree and contains the rows necessary
// to validate that the leaf belongs to the Tree. Leaves can be used by an
// Assembler to assemble a tree.
type Leaf struct {
	Data  []byte
	Rows  []ValidatorRow
	Index uint32
}

// Digest of the leaf, this can be checked against a known Tree hash before
// starting assembly. The digest of each additional Leaf will be checked by the
// assembler.
func (l *Leaf) Digest(h hash.Hash, buf []byte) []byte {
	buf = slice.NewBuffer(buf).Empty(cmpr.Max(4, h.Size()))

	h.Reset()
	rye.Serialize.Uint32(buf[:4], l.Index)
	h.Write(buf[:4])
	h.Write(l.Data)
	dig := h.Sum(buf[:0])
	for i := len(l.Rows) - 1; i >= 0; i-- {
		r := l.Rows[i]
		h.Reset()
		if r.IsFirst {
			h.Write(dig)
			h.Write(r.SiblingDigest)
		} else {
			h.Write(r.SiblingDigest)
			h.Write(dig)
		}
		dig = h.Sum(dig[:0])
	}
	return dig
}

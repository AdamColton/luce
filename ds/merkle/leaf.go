package merkle

import "hash"

// ValidatorRow represents a Branch in a Merkle Tree. It provides the digests
// needed for validation.
type ValidatorRow struct {
	Digests  [][]byte
	Child    int
	MaxIndex uint32
}

// Leaf represents a Data Leaf in a Merkle Tree and contains the rows necessary
// to validate that the leaf belongs to the Tree. Leaves can be used by an
// Assembler to assemble a tree.
type Leaf struct {
	Data  []byte
	Rows  []ValidatorRow
	Index uint32
}

func (l Leaf) branch(d int, h hash.Hash) (*branch, int) {
	row := l.Rows[d]
	r := len(row.Digests)
	if row.Child < 0 || row.Child > r {
		return nil, -1
	}
	b := &branch{
		children: make([]node, 0, r+1),
		idx:      row.MaxIndex,
	}

	descend := func() {
		if d+1 < len(l.Rows) {
			n, dr := l.branch(d+1, h)
			r += dr
			b.children = append(b.children, n)
		} else {
			b.children = append(b.children, newDataLeaf(l.Data, l.Index, h))
		}
	}

	for i, d := range row.Digests {
		if i == row.Child {
			descend()
		}
		b.children = append(b.children, digestNode(d))
	}
	if len(row.Digests) == row.Child {
		descend()
	}

	h.Reset()
	for _, c := range b.children {
		h.Write(c.Digest())
	}
	h.Write(uint32ToSlice(b.idx))
	b.digest = h.Sum(nil)

	return b, r
}

func (b *branch) Leaf(idx int) Leaf {
	return b.leaf(idx, Leaf{
		Rows: make([]ValidatorRow, 0, b.Depth()),
	})
}

func (b *branch) leaf(idx int, v Leaf) Leaf {
	r := ValidatorRow{
		Digests:  make([][]byte, 0, len(b.children)-1),
		MaxIndex: b.idx,
	}
	rIdx := len(v.Rows)
	v.Rows = append(v.Rows, r)
	for i, c := range b.children {
		ct := c.Count()
		if ct <= idx || idx == -1 {
			if idx != -1 {
				idx -= ct
			}
			r.Digests = append(r.Digests, c.Digest())
		} else {
			r.Child = i
			if idx == 0 {
				if dn, ok := c.(*dataLeaf); ok {
					v.Data = dn.data
					v.Index = dn.idx
				}
			}
			if nc, ok := c.(*branch); ok {
				v = nc.leaf(idx, v)
			}
			idx = -1
		}
	}
	v.Rows[rIdx] = r
	return v
}

// Digest of the leaf, this can be checked against a known Tree hash before
// starting assembly. The digest of each additional Leaf will be checked by the
// assembler.
func (l Leaf) Digest(h hash.Hash) []byte {
	return l.digest(0, h)
}

func (l Leaf) digest(d int, h hash.Hash) []byte {
	if d < 0 || d > len(l.Rows) {
		return nil
	}

	if d == len(l.Rows) {
		if l.Data == nil {
			return nil
		}
		h.Reset()
		h.Write(l.Data)
		h.Write(uint32ToSlice(l.Index))
		return h.Sum(nil)
	}

	cd := l.digest(d+1, h)
	row := l.Rows[d]
	h.Reset()
	for i, dig := range row.Digests {
		if i == row.Child {
			h.Write(cd)
		}
		h.Write(dig)
	}
	if row.Child == len(row.Digests) {
		h.Write(cd)
	}
	h.Write(uint32ToSlice(row.MaxIndex))
	return h.Sum(cd[0:0])
}

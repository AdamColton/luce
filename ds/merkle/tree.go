package merkle

import "hash"

type tree struct {
	node
	h             hash.Hash
	leaves, depth uint32
}

func (t *tree) Leaves() int {
	return int(t.leaves)
}

func (t *tree) Leaf(idx int) *Leaf {
	if idx < 0 || idx >= int(t.leaves) {
		return nil
	}
	l := &Leaf{
		Rows: make([]ValidatorRow, 0, t.depth),
	}
	n := t.node
	for {
		if dl, ok := n.(*dataLeaf); ok {
			l.Data = dl.data
			l.Index = dl.idx
			return l
		}
		b := n.(*branch)
		c0ls := b.children[0].Leaves()
		if idx < c0ls {
			n = b.children[0]
			l.Rows = append(l.Rows, ValidatorRow{
				SiblingDigest: b.children[1].Digest(),
				IsFirst:       true,
			})
		} else {
			n = b.children[1]
			idx -= c0ls
			l.Rows = append(l.Rows, ValidatorRow{
				SiblingDigest: b.children[0].Digest(),
				IsFirst:       false,
			})
		}
	}
}

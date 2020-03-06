package merkle

import (
	"bytes"
	"hash"
	"sort"
)

// Assembler of Validators. Adds Validators to a tree until the tree is
// complete.
type Assembler struct {
	root      *branch
	remaining int
	h         hash.Hash
}

// Assembler creates a new Assembler from the validator. All future validators
// added to the tree will be validated against against the initial validator.
func (l Leaf) Assembler(h hash.Hash) *Assembler {
	b, r := l.branch(0, h)
	return &Assembler{
		root:      b,
		remaining: r,
		h:         h,
	}
}

// Add a Leav to the tree being assembled. Returns a bool indicating if the
// Leaf was valid. Adding a Leaf multiple times will validate the Leaf, but
// will not change the tree. Add should not be called concurrently.
func (a *Assembler) Add(l Leaf) bool {
	return a.addLeaf(l, 0, a.root)
}

func (a *Assembler) addLeaf(l Leaf, depth int, b *branch) bool {
	if depth >= len(l.Rows) {
		return false
	}
	row := l.Rows[depth]

	// Validate
	if row.Child < 0 || row.Child > len(row.Digests) {
		return false
	}
	if !check(b.children[0:row.Child], row.Digests[0:row.Child]) ||
		!check(b.children[row.Child+1:], row.Digests[row.Child:]) {
		return false
	}

	c := b.children[row.Child]
	if b, ok := c.(*branch); ok {
		return a.addLeaf(l, depth+1, b)
	}
	if dig, ok := c.(digestNode); ok {
		if depth+1 == len(l.Rows) {
			dn := newDataLeaf(l.Data, l.Index, a.h)
			if !bytes.Equal(dig, dn.digest) {
				return false
			}
			b.children[row.Child] = dn
			a.remaining--
			return true
		}
		childB, r := l.branch(depth+1, a.h)
		if !bytes.Equal(childB.Digest(), dig) {
			return false
		}
		b.children[row.Child] = childB
		a.remaining += r - 1
		return true
	}
	if dl, ok := c.(*dataLeaf); ok {
		a.h.Reset()
		a.h.Write(l.Data)
		a.h.Write(uint32ToSlice(l.Index))
		return bytes.Equal(dl.digest, a.h.Sum(nil))
	}
	return false
}

// Done checks if the Tree is fully assembled. If it is the Tree is returned,
// if not, Tree will be null.
func (a *Assembler) Done() (bool, Tree) {
	if a.remaining == 0 {
		return true, a.root
	}
	return false, nil
}

func check(ns []node, d [][]byte) bool {
	if len(ns) != len(d) {
		return false
	}
	for i, n := range ns {
		if !bytes.Equal(n.Digest(), d[i]) {
			return false
		}
	}
	return true
}

// Need returns the indexes of all the Leaf Nodes that are still needed.
func (a *Assembler) Need() []uint32 {
	h := a.root.have(nil)
	sort.Slice(h, func(i, j int) bool {
		return h[i] < h[j]
	})
	var out []uint32
	for i := uint32(0); i <= a.root.idx; i++ {
		if len(h) > 0 && i == h[0] {
			h = h[1:]
		} else {
			out = append(out, i)
		}
	}
	return out
}

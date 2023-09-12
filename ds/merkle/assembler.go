package merkle

import (
	"bytes"
	"hash"

	"github.com/adamcolton/luce/serial/rye"
)

// Assembler of Leaves. Adds Leaves to a tree until the tree is
// complete.
type Assembler struct {
	root      node
	remaining int
	h         hash.Hash
	d         Description
	bools     *rye.Bits
}

// Done checks if assembly is done. If it is not, Tree will be nil.
func (a *Assembler) Done() (bool, Tree) {
	if a.remaining == 0 {
		t := &tree{
			node: a.root,
		}
		t.updateTree(a.h, true)
		return true, t
	}
	return false, nil
}

// Add a Leaf to the tree being assembled. Returns a bool indicating if the
// Leaf was valid. Adding a Leaf multiple times will validate the Leaf, but
// will not change the tree. Add should not be called concurrently.
func (a *Assembler) Add(l *Leaf) bool {
	ii := int(l.Index)
	if ii > a.bools.Ln {
		return false
	}
	a.bools.Idx = ii
	if a.bools.Read() > 0 || !bytes.Equal(a.d.Digest, l.Digest(a.h, nil)) {
		return false
	}

	a.bools.Idx = ii
	a.bools.Write(1)
	a.root = l.populate(0, a.root)
	a.remaining--

	return true
}

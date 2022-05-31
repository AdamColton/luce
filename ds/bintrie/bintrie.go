package bintrie

import "github.com/adamcolton/luce/serial/rye"

type node struct {
	branches [2]*node
	terminal bool
	size     int
}

func (n *node) private() {}
func (n *node) Size() int {
	return n.size
}

func (n *node) Insert(u uint32) {
	n.recurInsert(u, 32)
}

func (n *node) recurInsert(u, d uint32) {
	if d == 0 {
		n.terminal = true
	} else {
		n.getBranch(byte(u&1)).recurInsert(u>>1, d-1)
	}
	n.updateSize()
}

func (n *node) getSize() int {
	if n == nil {
		return 0
	}
	return n.size
}

func (n *node) updateSize() {
	n.size = n.branches[0].getSize() + n.branches[1].getSize()
	if n.terminal {
		n.size++
	}
}

func (n *node) getBranch(b byte) *node {
	if n.branches[b] == nil {
		n.branches[b] = &node{}
	}
	return n.branches[b]
}

func (n *node) Has(u uint32) bool {
	for i := 0; i < 32; i++ {
		b := u & 1
		u >>= 1
		if n.branches[b] == nil {
			return false
		}
		n = n.branches[b]
	}
	return n.terminal
}

func (n *node) All() []*rye.Bits {
	return n.recurAll(&rye.Bits{})
}

func (n *node) Copy() Trie {
	return n.copy()
}

func (n *node) copy() *node {
	if n == nil {
		return nil
	}
	return &node{
		branches: [2]*node{
			n.branches[0].copy(),
			n.branches[1].copy(),
		},
		terminal: n.terminal,
		size:     n.size,
	}
}

func (n *node) nilCheck() *node {
	if !n.terminal && n.branches[0] == nil && n.branches[1] == nil {
		return nil
	}
	return n
}
func boolCopy(bits *rye.Bits, bln bool) *rye.Bits {
	if bln {
		return bits.Copy()
	}
	return bits
}

func (n *node) recurAll(b *rye.Bits) (out []*rye.Bits) {
	if n.terminal {
		bt := boolCopy(b, n.branches[0] != nil || n.branches[1] != nil).Reset()
		out = append(out, bt)
	}
	if n.branches[0] != nil {
		b0 := boolCopy(b, n.branches[1] != nil)
		b0.Write(0)

		out = append(out, n.branches[0].recurAll(b0)...)
	}
	if n.branches[1] != nil {
		b.Write(1)
		out = append(out, n.branches[1].recurAll(b)...)
	}
	return
}

func Or(a, b Trie) Trie {
	return or(a.(*node), b.(*node))
}

func or(x, y *node) *node {
	t := x.terminal || y.terminal
	b0 := x.branches[0] != nil || y.branches[0] != nil
	b1 := x.branches[1] != nil || y.branches[1] != nil
	out := &node{
		terminal: t,
	}
	if b0 {
		if x.branches[0] == nil {
			out.branches[0] = y.branches[0].copy()
		} else if y.branches[0] == nil {
			out.branches[0] = x.branches[0].copy()
		} else {
			out.branches[0] = or(x.branches[0], y.branches[0])
		}
	}
	if b1 {
		if x.branches[1] == nil {
			out.branches[1] = y.branches[1].copy()
		} else if y.branches[1] == nil {
			out.branches[1] = x.branches[1].copy()
		} else {
			out.branches[1] = or(x.branches[1], y.branches[1])
		}
	}
	out.updateSize()
	return out
}

func And(a, b Trie) Trie {
	return and(a.(*node), b.(*node))
}

func and(x, y *node) *node {
	t := x.terminal && y.terminal
	and0 := x.branches[0] != nil && y.branches[0] != nil
	and1 := x.branches[1] != nil && y.branches[1] != nil
	if !t && !and0 && !and1 {
		return nil
	}
	n := &node{
		terminal: t,
	}
	if and0 {
		n.branches[0] = and(x.branches[0], y.branches[0])
	}
	if and1 {
		n.branches[1] = and(x.branches[1], y.branches[1])
	}
	n.updateSize()
	return n.nilCheck()
}

// Nand returns all values in a but not in b.
func Nand(a, b Trie) Trie {
	return nand(a.(*node), b.(*node))
}

func nand(x, y *node) *node {
	t := x.terminal && !y.terminal
	x0 := x.branches[0] != nil
	x1 := x.branches[1] != nil
	if !t && !x0 && !x1 {
		return nil
	}
	n := &node{}
	if x0 {
		if y.branches[0] != nil {
			n.branches[0] = nand(x.branches[0], y.branches[0])
		} else {
			n.branches[0] = x.branches[0].copy()
		}
	}
	if x1 {
		if y.branches[1] != nil {
			n.branches[1] = nand(x.branches[1], y.branches[1])
		} else {
			n.branches[1] = x.branches[1].copy()
		}
	}
	n.updateSize()
	return n.nilCheck()
}

func (n *node) Delete(u uint32) {
	n.recurDelete(u, 32)
}

func (n *node) recurDelete(u, d uint32) *node {
	if n == nil {
		return nil
	}
	if d == 0 {
		n.terminal = false
	} else {
		b := u & 1
		n.branches[b] = n.branches[b].recurDelete(u>>1, d-1)
	}
	n.updateSize()
	return n.nilCheck()
}

func (n *node) InsertTrie(t Trie) {
	n.recurInsertTrie(t.(*node))
}

func (n *node) recurInsertTrie(n2 *node) {
	b0 := n2.branches[0] != nil
	b1 := n2.branches[1] != nil

	if b0 {
		if n.branches[0] == nil {
			n.branches[0] = n2.branches[0].copy()
		} else {
			n.branches[0].recurInsertTrie(n2.branches[0])
		}
	}
	if b1 {
		if n.branches[1] == nil {
			n.branches[1] = n2.branches[1].copy()
		} else {
			n.branches[1].recurInsertTrie(n2.branches[1])
		}
	}
	n.updateSize()
}

package bintrie

import "github.com/adamcolton/luce/serial/rye"

type node[U Uint] struct {
	branches [2]*node[U]
	terminal bool
	size     int
}

func (n *node[U]) private() {}
func (n *node[U]) Size() int {
	return n.size
}

func (n *node[U]) Insert(u U) {
	n.recurInsert(u, sizeOf(u))
}

func (n *node[U]) recurInsert(u, d U) {
	if d == 0 {
		n.terminal = true
	} else {
		n.getBranch(byte(u&1)).recurInsert(u>>1, d-1)
	}
	n.updateSize()
}

func (n *node[U]) getSize() int {
	if n == nil {
		return 0
	}
	return n.size
}

func (n *node[U]) updateSize() {
	n.size = n.branches[0].getSize() + n.branches[1].getSize()
	if n.terminal {
		n.size++
	}
}

func (n *node[U]) getBranch(b byte) *node[U] {
	if n.branches[b] == nil {
		n.branches[b] = &node[U]{}
	}
	return n.branches[b]
}

func (n *node[U]) Has(u U) bool {
	for i := U(0); i < sizeOf(u); i++ {
		b := u & 1
		u >>= 1
		if n.branches[b] == nil {
			return false
		}
		n = n.branches[b]
	}
	return n.terminal
}

func (n *node[U]) All() []*rye.Bits {
	return n.recurAll(&rye.Bits{})
}

func (n *node[U]) Copy() Trie[U] {
	return n.copy()
}

func (n *node[U]) copy() *node[U] {
	if n == nil {
		return nil
	}
	return &node[U]{
		branches: [2]*node[U]{
			n.branches[0].copy(),
			n.branches[1].copy(),
		},
		terminal: n.terminal,
		size:     n.size,
	}
}

func (n *node[U]) nilCheck() *node[U] {
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

func (n *node[U]) recurAll(b *rye.Bits) (out []*rye.Bits) {
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

func (n *node[U]) Delete(u U) {
	n.recurDelete(u, 32)
}

func (n *node[U]) recurDelete(u, d U) *node[U] {
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

func (n *node[U]) InsertTrie(t Trie[U]) {
	n.recurInsertTrie(t.(*node[U]))
}

func (n *node[U]) recurInsertTrie(n2 *node[U]) {
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

func (n *node[U]) Union(t Trie[U]) {
	n.recurUnion(t.(*node[U]))
}

func (n *node[U]) recurUnion(n2 *node[U]) *node[U] {
	t := n.terminal && n2.terminal
	and0 := n.branches[0] != nil && n2.branches[0] != nil
	and1 := n.branches[1] != nil && n2.branches[1] != nil
	if !t && !and0 && !and1 {
		return nil
	}
	n.terminal = t
	if and0 {
		n.branches[0] = n.branches[0].recurUnion(n2.branches[0])
	} else {
		n.branches[0] = nil
	}
	if and1 {
		n.branches[1] = n.branches[1].recurUnion(n2.branches[1])
	} else {
		n.branches[1] = nil
	}
	n.updateSize()
	return n.nilCheck()
}

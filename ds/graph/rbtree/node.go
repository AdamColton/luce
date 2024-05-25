package rbtree

import (
	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/serial/rye"
)

type color int8

const (
	red   color = -1
	black color = 1

	left  = 0
	right = 1
)

type Node[Key, Val any] struct {
	graph.KV[Key, Val]
	chld [2]graph.Ptr[*Node[Key, Val]]
	prt  graph.Ptr[*Node[Key, Val]]
	color
	size int
	id   uint32
}

func (n *Node[Key, Val]) EntKey() []byte {
	if n == nil {
		return nil
	}
	out := make([]byte, 4)
	rye.Serialize.Uint32(out, n.id)
	return out
}

func stripBool[T any](t T, b bool) T {
	return t
}

func (n *Node[Key, Val]) AtIdx(idx int) graph.Node[Key, Val] {
	if idx >= 0 && idx < 2 {
		return stripBool(n.chld[idx].Get())
	}
	return nil
}

func (n *Node[Key, Val]) Len() int {
	return 2
}

func (n *Node[Key, Val]) getSize() int {
	if n == nil {
		return 0
	}
	return n.size
}

func (n *Node[Key, Val]) setChld(idx int, c *Node[Key, Val]) {
	if n != nil {
		n.chld[idx] = n.chld[idx].Set(c)
	}
	c.setPrt(n)
}

func (n *Node[Key, Val]) setPrt(p *Node[Key, Val]) {
	// TODO: check for redundant calls (set child already calls this)
	if n != nil {
		n.prt = n.prt.Set(p)
	}
}

func (n *Node[Key, Val]) rot(idx int) {
	if n == nil {
		return
	}
	idxi := idxInv(idx)

	p := stripBool(n.prt.Get())
	pidx := n.idx()
	c := n.getChild(idxi)
	trade := c.getChild(idx)
	checkLoop(n)
	//c.setPrt(c.prt.Get())
	//n.setPrt(c)
	n.setChld(idxi, trade)
	c.setChld(idx, n)
	p.setChld(pidx, c)
	checkLoop(n)

	n.size = n.getChild(0).getSize() + n.getChild(1).getSize() + 1
	c.size = c.getChild(0).getSize() + c.getChild(1).getSize() + 1
}

func (n *Node[Key, Val]) clr() color {
	if n == nil {
		return black
	}
	return n.color
}

func (n *Node[Key, Val]) set(c color) {
	if n == nil {
		return
	}
	n.color = c
}

func (n *Node[Key, Val]) uncle(prt, gprt *Node[Key, Val]) *Node[Key, Val] {
	if gprt == nil {
		return nil
	}
	c0 := gprt.getChild(0)
	if c0 == prt {
		return gprt.getChild(1)
	}
	return c0
}

func (n *Node[Key, Val]) ancestors() (prt *Node[Key, Val], gprt *Node[Key, Val]) {
	prt, _ = n.prt.Get()
	if prt != nil {
		gprt, _ = prt.prt.Get()
	}
	return
}

func (n *Node[Key, Val]) idx() (i int) {
	p, _ := n.prt.Get()
	if p != nil && p.getChild(i) != n {
		i = right
	}
	return
}

// next node in the given direction, so the value immediatly larger or smaller.
func (n *Node[Key, Val]) next(dir int) *Node[Key, Val] {
	c := n
	for c.prt != nil && stripBool(c.prt.Get()).getChild(dir) == c {
		nxt := stripBool(c.prt.Get()).getChild(dir)
		if nxt == nil || nxt.getChild(dir) == c {
			break
		}
		c = stripBool(c.prt.Get())
	}
	if c.prt == nil {
		return nil
	}

	return stripBool(c.prt.Get()).leaf(dir)
}

func (n *Node[Key, Val]) leaf(dir int) *Node[Key, Val] {
	for {
		next := n.getChild(dir)
		if next == nil {
			break
		}
		n = next
	}
	return n
}

func (n *Node[Key, Val]) sibling() (sib *Node[Key, Val]) {
	p := stripBool(n.prt.Get())
	if p == nil {
		return
	}
	sib = p.getChild(left)
	if sib == n {
		sib = p.getChild(right)
	}
	return
}

func (n *Node[Key, Val]) getChild(idx int) *Node[Key, Val] {
	if n == nil {
		return nil
	}
	c, _ := n.chld[idx].Get()
	return c
}

// remove n and promote child cidx. Expect that the other child is nil.
func (n *Node[Key, Val]) remove(cidx int) {
	c := n.getChild(cidx)
	p := stripBool(n.prt.Get())
	c.setPrt(p)
	if p != nil {
		p.setChld(n.idx(), c)
	}
}

func idxInv(idx int) int {
	return 1 - idx
}

func checkLoop[Key, Val any](n *Node[Key, Val]) {
	check := lset.New(n)
	for {
		next := stripBool(n.prt.Get())
		if next == nil {
			return
		}
		if check.Contains(next) {
			panic("loop")
		}
		check.Add(next)
		n = next
	}
}

func fixupInsert[Key, Val any](n *Node[Key, Val]) *Node[Key, Val] {
	prt, gprt := n.ancestors()
	for {
		if prt.clr() == black {
			return n
		}

		uncle := n.uncle(prt, gprt)
		if uncle.clr() == black {
			break
		}
		prt.set(black)
		uncle.set(black)
		gprt.set(red)
		n = gprt
		prt, gprt = n.ancestors()
	}

	idx := n.idx()
	idxi := idxInv(idx)

	if prt == gprt.getChild(idxi) {
		prt.rot(idxi)
		checkLoop(n)
		prt, gprt = n, stripBool(n.prt.Get())
		idx, idxi = idxi, idx
	}

	prt.set(black)
	gprt.set(red)
	if prt == gprt.getChild(idx) {
		gprt.rot(idxi)
		checkLoop(n)
	}
	return n
}

func fixupDelete[Key, Val any](n *Node[Key, Val]) {
	var idx int
	sib := n.sibling()
	prt := stripBool(n.prt.Get())
	for {
		if prt == nil {
			return
		}
		idx = n.idx()

		if sib.clr() == red {
			prt.set(red)
			sib.set(black)
			prt.rot(idx)
			sib = n.sibling()
		}

		cond := sib.clr() == black &&
			sib.getChild(left).clr() == black &&
			sib.getChild(right).clr() == black

		if !cond {
			break
		}

		if prt.clr() == red {
			sib.set(red)
			prt.set(black)
			return
		}

		sib.set(red)
		n = prt
		prt = stripBool(n.prt.Get())
		sib = n.sibling()
	}

	idxi := idxInv(idx)

	cond := sib.clr() == black &&
		sib.getChild(idx).clr() == red &&
		sib.getChild(idxi).clr() == black

	if cond {
		sib.set(red)
		sib.getChild(idx).set(black)
		sib.rot(idxi)
		sib = n.sibling()
	}

	sib.set(prt.clr())
	prt.set(black)
	if sii := sib.getChild(idxi); sii.clr() == red {
		sii.set(black)
		prt.rot(idx)
	}
}

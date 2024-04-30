package rbtree

import (
	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/util/filter"
)

type Tree[Key, Val any] struct {
	root graph.Ptr[*node[Key, Val]]
	cmpr filter.Compare[Key]
	size int
}

func MakePtrType[Key, Val any]() graph.Ptr[*node[Key, Val]] {
	return graph.RawPointer[node[Key, Val]]{}
}

func New[Key, Val any](ptr graph.Ptr[*node[Key, Val]], cmpr filter.Compare[Key]) *Tree[Key, Val] {
	return &Tree[Key, Val]{
		cmpr: cmpr,
		root: ptr,
	}
}

func rcmpr(x int) int {
	return (3*x*x + x - 2) / 2
}

func (t *Tree[Key, Val]) Root() graph.Node[Key, Val] {
	return t.root.Get()
}

func (t *Tree[Key, Val]) Add(k Key, v Val) {
	add := &node[Key, Val]{
		KV:    graph.NewKV(k, v),
		color: red,
		size:  1,
		chld: [2]graph.Ptr[*node[Key, Val]]{
			t.root.New(),
			t.root.New(),
		},
		prt: t.root.New(),
	}
	r := t.root.Get()
	if r == nil {
		add.set(black)
		t.root = t.root.Set(add)
	} else {
		n := r
		for {
			n.size++
			idx := rcmpr(t.cmpr(k, n.KV.K))
			if idx < 0 {
				break
			}
			nxt := n.getChild(idx)
			if nxt == nil {
				n.setChld(idx, add)
				n = add
				t.fixRoot(fixupInsert(n))
				break
			}
			n = nxt
		}
		n.KV = graph.NewKV(k, v)
	}

	t.size++
}

func (t *Tree[Key, Val]) Remove(k Key) {
	n, found := t.seek(k)
	if !found {
		return
	}
	c0, c1 := n.getChild(0), n.getChild(1)
	if c0 != nil && c1 != nil {
		pred := c0.leaf(1)
		n.KV = pred.KV
		n = pred
		c0, c1 = n.getChild(0), n.getChild(1)
	}
	if c0 == nil || c1 == nil {
		cidx := left
		c := c0
		if c == nil {
			cidx = right
			c = c1
		}

		prt := n.prt.Get()
		fixup := prt
		if fixup == nil {
			fixup = c
		}

		if n.clr() == black {
			n.set(c.clr())
			fixupDelete(n)
		}
		if prt != nil {
			dec := prt
			for {
				dec.size--
				nxt := dec.prt.Get()
				if nxt == nil {
					break
				}
				dec = nxt
			}
		}
		n.remove(cidx)
		if fixup == nil {
			t.root = t.root.Set(nil)
		} else {
			t.fixRoot(fixup)
		}
	}
	t.size--
}

func (t *Tree[Key, Val]) fixRoot(n *node[Key, Val]) {
	check := lset.New(n)
	for {
		next := n.prt.Get()
		if next == nil {
			break
		}
		if check.Contains(next) {
			panic("loop")
		}
		check.Add(next)
		n = next
	}
	n.set(black)
	t.root = t.root.Set(n)

}

func (t *Tree[Key, Val]) Seek(k Key) (n Node[Key, Val], found bool) {
	n, found = t.seek(k)
	return
}

func (t *Tree[Key, Val]) seek(k Key) (n *node[Key, Val], found bool) {
	n = t.root.Get()
	if n == nil {
		return
	}

	var idx int
	for {
		idx = rcmpr(t.cmpr(k, n.K))
		if idx < 0 {
			break
		}
		nxt := n.getChild(idx)
		if nxt == nil {
			break
		}
		n = nxt
	}
	found = idx == -1
	if !found && idx == 1 {
		n = n.next(1)
	}
	return
}

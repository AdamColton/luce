package ll

import (
	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/util/liter"
)

type Double[Key, Val any] struct {
	graph.KV[Key, Val]
	Prev, Next *Double[Key, Val]
}

func NewDouble[Key, Val any](k Key, v Val) *Double[Key, Val] {
	return &Double[Key, Val]{
		KV: graph.NewKV(k, v),
	}
}

func NewDoubleLoop[Key, Val any](k Key, v Val) *Double[Key, Val] {
	n := NewDouble(k, v)
	n.Prev = n
	n.Next = n
	return n
}

func (d *Double[Key, Val]) InsertAfter(k Key, v Val) *Double[Key, Val] {
	n := &Double[Key, Val]{
		KV:   graph.NewKV(k, v),
		Next: d.Next,
		Prev: d,
	}
	d.Next = n
	if n.Next != nil {
		n.Next.Prev = n
	}
	return n
}

func (d *Double[Key, Val]) Remove() {
	if d.Prev != nil {
		d.Prev.Next = d.Next
	}
	if d.Next != nil {
		d.Next.Prev = d.Prev
	}
	d.Next = nil
	d.Prev = nil
}

// AtIdx return Next if given 0 and Prev if given 1.
func (d *Double[Key, Val]) AtIdx(idx int) graph.Node[Key, Val] {
	if idx == 0 {
		return d.Next
	} else if idx == 1 {
		return d.Prev
	}
	return nil
}

func (s *Double[Key, Val]) Len() int {
	return 2
}

func (d *Double[Key, Val]) Iter(forward bool) liter.Iter[graph.KV[Key, Val]] {
	return &doubleIter[Key, Val]{
		n:       d,
		start:   d,
		forward: forward,
	}
}

func (d *Double[Key, Val]) DirFactory(forward bool) (iter liter.Iter[graph.KV[Key, Val]], kv graph.KV[Key, Val], done bool) {
	done = d == nil
	if !done {
		kv = d.KV
	}
	iter = d.Iter(forward)
	return
}

func (d *Double[Key, Val]) ForwardFactory() (iter liter.Iter[graph.KV[Key, Val]], kv graph.KV[Key, Val], done bool) {
	return d.DirFactory(true)
}

func (d *Double[Key, Val]) BackwardFactory() (iter liter.Iter[graph.KV[Key, Val]], kv graph.KV[Key, Val], done bool) {
	return d.DirFactory(false)
}

type doubleIter[Key, Val any] struct {
	i        int
	start, n *Double[Key, Val]
	forward  bool
}

func (di *doubleIter[Key, Val]) Next() (kv graph.KV[Key, Val], done bool) {
	done = di.Done()
	if done {
		return
	}
	di.i++
	if di.forward {
		di.n = di.n.Next
	} else {
		di.n = di.n.Prev
	}
	done = di.Done()
	if !done {
		kv = di.n.KV
	}
	return
}

func (di *doubleIter[Key, Val]) Cur() (kv graph.KV[Key, Val], done bool) {
	done = di.Done()
	if !done {
		kv = di.n.KV
	}
	return
}

func (si *doubleIter[Key, Val]) Done() bool {
	return si.n == nil || (si.i > 0 && si.n == si.start)
}

func (di *doubleIter[Key, Val]) Idx() int {
	return di.i
}

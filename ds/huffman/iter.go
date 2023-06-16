package huffman

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/util/iter"
)

// HuffIter fulfills iter.Iter as well as iter.Starter and has a Factory method
// to make a copy and fulfill iter.Factory.
type HuffIter[T any] interface {
	iter.Iter[T]
	iter.Starter[T]
	Factory() (copy iter.Iter[T], t T, done bool)
	slice.Slicer[T]
}

type huffiter[T any] struct {
	node *huffNode[T]
	b    *rye.Bits
	t    T
	idx  int
}

func (i *huffiter[T]) Next() (t T, done bool) {
	done = i.Done()
	if !done {
		i.t = i.node.Read(i.b)
		t = i.t
		i.idx++
	}
	return
}

func (i *huffiter[T]) Cur() (t T, done bool) {
	done = i.Done()
	if !done {
		t = i.t
	}
	return
}

func (i *huffiter[T]) Done() bool {
	return i.b.Idx >= i.b.Ln
}

func (i *huffiter[T]) Idx() int {
	return i.idx
}

func (i *huffiter[T]) Start() (t T, done bool) {
	i.idx = -1
	i.b.Idx = 0
	return i.Next()
}

func (i *huffiter[T]) Factory() (copy iter.Iter[T], t T, done bool) {
	copy = &huffiter[T]{
		node: i.node,
		b:    i.b.ShallowCopy().Reset(),
		idx:  -1,
	}
	t, done = copy.Next()
	return
}

func (i *huffiter[T]) Slice(buf []T) []T {
	cp := i.b.ShallowCopy().Reset()
	out := buf
	for cp.Idx < cp.Ln {
		out = append(out, i.node.Read(cp))
	}
	return out
}

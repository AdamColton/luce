package huffman

import (
	"sort"

	"github.com/adamcolton/luce/ds/heap"
	"github.com/adamcolton/luce/serial/rye"
)

// Tree represents a Huffman Coding. Create one with New or MapNew. It has an
// unexported method, so it can only be fulfilled in this package.
type Tree[T any] interface {
	// Read decodes one value from b, moving b.Idx past its bits.
	Read(b *rye.Bits) T
	// Iter returns an iterator that decodes all of the bits in b, from the
	// beginning of b. The iterator shares b, so it moves b.Idx.
	Iter(b *rye.Bits) HuffIter[T]
	// All visits every value in the Tree can calls the given func on each value,
	// from the leftmost branch (the 0 bit) to the rightmost.
	All(fn func(T))
	// Len returns the number of values in the Tree.
	Len() int
	private()
}

type tree[T any] struct {
	ln int
	*huffNode[T]
}

func (tr tree[T]) Len() int {
	return tr.ln
}

func (tr tree[T]) Iter(b *rye.Bits) HuffIter[T] {
	i := &huffiter[T]{
		node: tr.huffNode,
		b:    b,
	}
	i.Start()

	return i
}

func (tr tree[T]) private() {}

// Frequency is used for constructing a Huffman Coding.
type Frequency[T any] struct {
	Val   T
	Count int
}

// New Huffman Coding Tree constructed from Frequency data. It panics if data is
// empty. Each element is a separate leaf, so a value that appears more than once
// in data is in the Tree more than once.
func New[T any](data []Frequency[T]) Tree[T] {
	h := newHeap[T](len(data))
	for _, d := range data {
		h.Data = append(h.Data, d.root())
	}
	sort.Slice(h.Data, h.Less)

	return tree[T]{
		huffNode: makeHeapTree(h),
		ln:       len(data),
	}
}

// MapNew creates a Huffman Coding Tree constructed from a frequency map. It
// panics if data is empty.
func MapNew[T comparable](data map[T]int) Tree[T] {
	h := newHeap[T](len(data))
	for v, c := range data {
		h.Data = append(h.Data, newLeaf(v, c))
	}
	sort.Slice(h.Data, h.Less)

	return tree[T]{
		huffNode: makeHeapTree(h),
		ln:       len(data),
	}
}

func newHeap[T any](ln int) *heap.Heap[*root[T]] {
	h := &heap.Heap[*root[T]]{
		Data: make([]*root[T], 0, ln),
	}
	h.Less = func(i, j int) bool {
		return h.Data[i].sum < h.Data[j].sum
	}
	return h
}

func makeHeapTree[T any](data *heap.Heap[*root[T]]) *huffNode[T] {
	for ln := len(data.Data); ln > 1; ln-- {
		a, b := data.Pop(), data.Pop()
		data.Push(newBranch(a.node, b.node, a.sum+b.sum))
	}
	return data.Data[0].node
}

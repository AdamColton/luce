package huffman

import (
	"sort"

	"github.com/adamcolton/luce/ds/heap"
	"github.com/adamcolton/luce/serial/rye"
)

// Tree represents a Huffman Coding.
type Tree[T any] interface {
	Read(b *rye.Bits) T
	ReadAll(b *rye.Bits) []T
}

// Frequency is used for constructing a Huffman Coding.
type Frequency[T any] struct {
	Val   T
	Count int
}

// New Huffman Coding Tree contructed from Frequency data.
func New[T any](data []Frequency[T]) Tree[T] {
	h := newHeap[T](len(data))
	for _, d := range data {
		h.Data = append(h.Data, d.root())
	}
	sort.Slice(h.Data, h.Less)

	return makeHeapTree(h)
}

// New Huffman Coding Tree contructed from a frequency map.
func MapNew[T comparable](data map[T]int) Tree[T] {
	h := newHeap[T](len(data))
	for v, c := range data {
		h.Data = append(h.Data, newLeaf(v, c))
	}
	sort.Slice(h.Data, h.Less)

	return makeHeapTree(h)
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

func makeHeapTree[T any](data *heap.Heap[*root[T]]) Tree[T] {
	for ln := len(data.Data); ln > 1; ln-- {
		a, b := data.Pop(), data.Pop()
		data.Push(newBranch(a.node, b.node, a.sum+b.sum))
	}
	return data.Data[0].node
}

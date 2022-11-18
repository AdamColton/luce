package heap

import "sort"

// Heap keeps data sorted with O(log n) for Push and Pop.
type Heap[T any] struct {
	Data []T
	Less func(i, j int) bool
}

// Ordered data types can be used with NewMin and NewMax without the need to
// define the Less func on the Heap.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 |
		~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

// NewMin creates a min heap. Pop will return the lowest value.
func NewMin[T Ordered]() *Heap[T] {
	h := &Heap[T]{}
	h.Less = func(i, j int) bool {
		return h.Data[i] < h.Data[j]
	}
	return h
}

// NewMax creates a max heap. Pop will return the highest value.
func NewMax[T Ordered]() *Heap[T] {
	h := &Heap[T]{}
	h.Less = func(i, j int) bool {
		return h.Data[i] > h.Data[j]
	}
	return h
}

// SetLess using a function to compare the types directly instead of comparing
// based in indice.
func (h *Heap[T]) SetLess(less func(i, j T) bool) {
	h.Less = func(i, j int) bool {
		return less(h.Data[i], h.Data[j])
	}
}

// Push a value onto the heap.
func (h *Heap[T]) Push(v T) {
	h.Data = append(h.Data, v)
	h.bubbleUp(len(h.Data) - 1)
}

// Pop a value from the heap. This will return the "Least" value according to
// the Less func.
func (h *Heap[T]) Pop() T {
	out := h.Data[0]
	ln := len(h.Data) - 1
	h.Data[0] = h.Data[ln]
	h.Data = h.Data[:ln]
	h.bubbleDown(0)
	return out
}

// Sort a heap using Less. Useful when initilizing a heap.
func (h *Heap[T]) Sort() {
	sort.Slice(h.Data, h.Less)
}

func (h *Heap[T]) swap(i, j int) {
	h.Data[i], h.Data[j] = h.Data[j], h.Data[i]
}

func parent(idx int) int {
	return (idx - 1) / 2
}

func left(idx int) int {
	return idx*2 + 1
}

func right(idx int) int {
	return idx*2 + 2
}

func (h *Heap[T]) bubbleUp(idx int) {
	for pIdx := parent(idx); pIdx >= 0 && h.Less(idx, pIdx); pIdx, idx = parent(pIdx), pIdx {
		h.swap(idx, pIdx)
	}
}

func (h *Heap[T]) bubbleDown(idx int) {
	ln := len(h.Data)
	for l, r := left(idx), right(idx); l < ln; idx, l, r = l, left(l), right(l) {
		if r < ln && h.Less(r, l) {
			l = r
		}
		if h.Less(l, idx) {
			h.swap(l, idx)
		}
	}
}

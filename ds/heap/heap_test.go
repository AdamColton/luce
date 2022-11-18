package heap

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeap(t *testing.T) {
	vals := make([]int, 100)
	for i := range vals {
		vals[i] = rand.Intn(1000)
	}
	m := NewMin[int]()
	for _, v := range vals {
		m.Push(v)
	}
	M := NewMax[int]()
	M.Data = make([]int, len(vals))
	copy(M.Data, vals)
	M.Sort()

	h := &Heap[int]{
		Data: make([]int, len(vals)),
	}
	copy(h.Data, vals)

	sort.Slice(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})

	for _, v := range vals {
		assert.Equal(t, v, m.Pop())
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i] > vals[j]
	})

	for _, v := range vals {
		assert.Equal(t, v, M.Pop())
	}

	less := func(i, j int) bool {
		if i%2 == 0 && j%2 == 1 {
			return true
		} else if i%2 == 1 && j%2 == 0 {
			return false
		}
		return i < j
	}
	h.SetLess(less)
	h.Sort()

	sort.Slice(vals, func(i, j int) bool {
		return less(vals[i], vals[j])
	})
	for _, v := range vals {
		assert.Equal(t, v, h.Pop())
	}
}

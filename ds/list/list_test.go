package list_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/util/iter"
	"github.com/stretchr/testify/assert"
)

func TestLists(t *testing.T) {
	pi := []int{3, 1, 4, 1, 5}
	tt := map[string]struct {
		expected []int
		list.List[int]
	}{
		"SliceList": {
			expected: pi,
			List:     list.SliceList[int](pi),
		},
		"Generator": {
			expected: []int{0, 1, 4, 9, 16},
			List: list.Generator[int]{
				Fn: func(i int) int {
					return i * i
				},
				Length: 5,
			},
		},
		"Reverse": {
			expected: []int{5, 1, 4, 1, 3},
			List:     list.Reverse[int]{list.SliceList[int]([]int{3, 1, 4, 1, 5})},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, len(tc.expected), tc.List.Len())

			got := list.ToSlice(tc.List, nil)
			assert.Equal(t, tc.expected, got)

			f := list.IterFactory(tc.List)
			var last int
			fn := func(i, idx int) {
				assert.Equal(t, tc.expected[idx], i)
				last = idx
			}
			f.For(fn)
			assert.Len(t, tc.expected, last+1)

			it := list.NewIter(tc.List)
			last = 0
			iter.For[int](it, fn)
			assert.Len(t, tc.expected, last+1)

			it.Start()
			assert.Equal(t, 0, it.Idx())
			got = got[:0]
			iter.For[int](it, func(t int, idx int) {
				got = append(got, t)
			})
			assert.Equal(t, tc.expected, got)

			_, done := it.Next()
			assert.True(t, done)
		})
	}
}

func ExampleList_toSlice() {
	l := list.Generator[int]{
		Fn: func(i int) int {
			return i * i
		},
		Length: 5,
	}
	// To convert a List to a slice, create an Iter and use slice.IterSlice
	// If the underlying List fulfills Slicer, that will be invoked.
	// Otherwise, IterSlice will iterate over the list to create a slice.
	s := list.ToSlice[int](l, nil)
	fmt.Println(s)
	// Output:
	// [0 1 4 9 16]
}

func TestFoo(t *testing.T) {
	u := uint16(55)
	s := int(unsafe.Sizeof(u)) * 8
	assert.Equal(t, 32, s)
}

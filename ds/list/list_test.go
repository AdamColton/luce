package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
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
			s := list.ToSlice(tc.List, nil)
			assert.Equal(t, tc.expected, s)
			assert.Equal(t, len(tc.expected), tc.List.Len())
			list.NewIter(tc.List).Do(func(idx, i int) bool {
				assert.Equal(t, tc.expected[idx], i)
				return false
			})

		})
	}
}

package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/stretchr/testify/assert"
)

func TestLists(t *testing.T) {
	tt := map[string]struct {
		expected []int
		list.List[int]
	}{
		"Generator": {
			expected: []int{0, 1, 4, 9, 16},
			List: list.Generator[int]{
				Fn: func(i int) int {
					return i * i
				},
				Length: 5,
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, len(tc.expected), tc.List.Len())
			for i, e := range tc.expected {
				assert.Equal(t, e, tc.List.AtIdx(i))
			}
		})
	}
}

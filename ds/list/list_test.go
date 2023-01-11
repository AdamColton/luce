package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/util/iter"
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

			it.I = it.Len()
			_, done := it.Next()
			assert.True(t, done)
		})
	}
}

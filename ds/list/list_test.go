package list_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

type mockList struct{}

func (m mockList) AtIdx(idx int) int { return idx }
func (m mockList) Len() int          { return 10 }
func (m mockList) String() string    { return "0...9" }

func TestWrap(t *testing.T) {
	var m mockList
	w := list.Wrap[int](m)
	assert.Equal(t, m, w.List)

	w = list.Wrap[int](w)
	assert.Equal(t, m, w.List)
	_, shouldBeFalse := w.List.(list.Wrapper[int])
	assert.False(t, shouldBeFalse)
}

func TestUpgrade(t *testing.T) {
	var m mockList
	w := list.Wrap[int](m)
	var s fmt.Stringer
	assert.True(t, upgrade.Upgrade(w, &s))
	assert.Equal(t, "0...9", s.String())

	s = nil
	i := w.Iter()
	assert.True(t, upgrade.Upgrade(i, &s))
	assert.Equal(t, "0...9", s.String())
}

func TestLists(t *testing.T) {
	pi := []int{3, 1, 4, 1, 5}
	tt := map[string]struct {
		expected []int
		list.Wrapper[int]
	}{
		"SliceList": {
			expected: pi,
			Wrapper:  list.Slice(pi),
		},
		"Generator": {
			expected: []int{0, 1, 4, 9, 16},
			Wrapper: list.Generator[int]{
				Fn: func(i int) int {
					return i * i
				},
				Length: 5,
			}.Wrap(),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			w := tc.Wrapper
			assert.Equal(t, len(tc.expected), w.Len())

			fn := func(i, idx int) {
				assert.Equal(t, tc.expected[idx], i)
			}

			c := w.IterFactory().ForIdx(fn)
			assert.Len(t, tc.expected, c)

			it := w.Iter()
			var s iter.Starter[int]
			assert.True(t, upgrade.Upgrade(it, &s))

			c = it.ForIdx(fn)
			assert.Len(t, tc.expected, c)

			s.Start()
			assert.Equal(t, 0, it.Idx())

			var ls list.Slicer[int]
			if upgrade.Upgrade(w, &ls) {
				assert.Equal(t, tc.expected, ls.Slice(nil))
			}

			got := iter.Appender[int]().
				Iter(nil, it)
			assert.Equal(t, tc.expected, got)

			_, done := it.Next()
			assert.True(t, done)
		})
	}
}

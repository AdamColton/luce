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
	s, ok := upgrade.To[fmt.Stringer](w)
	assert.True(t, ok)
	assert.Equal(t, "0...9", s.String())

	i := w.Iter()
	s, ok = upgrade.To[fmt.Stringer](i)
	assert.True(t, ok)
	assert.Equal(t, "0...9", s.String())
}

func iSq(i int) int {
	return i * i
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
			Wrapper:  list.GeneratorFactory(iSq)(5).Wrap(),
		},
		"Reverse": {
			expected: []int{5, 1, 4, 1, 3},
			Wrapper:  list.Slice(pi).Reverse(),
		},
		"Transformer": {
			expected: []int{20, 5, 19, 20, 9, 14, 7},
			Wrapper: list.Transformer[rune, int]{
				List: list.Slice([]rune("testing")),
				Fn: func(r rune) int {
					return int(r) - int('a') + 1
				},
			}.Wrap(),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			w := tc.Wrapper
			assert.Equal(t, len(tc.expected), w.Len())

			got := w.ToSlice(nil)
			assert.Equal(t, tc.expected, got)

			fn := func(i, idx int) {
				assert.Equal(t, tc.expected[idx], i)
			}

			c := w.IterFactory().ForIdx(fn)
			assert.Len(t, tc.expected, c)

			it := w.Iter()
			s, ok := upgrade.To[iter.Starter[int]](it)
			assert.True(t, ok)

			c = it.ForIdx(fn)
			assert.Len(t, tc.expected, c)

			s.Start()
			assert.Equal(t, 0, it.Idx())

			got = iter.Appender[int]().
				Iter(got[:0], it)
			assert.Equal(t, tc.expected, got)

			_, done := it.Next()
			assert.True(t, done)
		})
	}
}

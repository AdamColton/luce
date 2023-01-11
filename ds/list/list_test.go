package list_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/list"
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
}

func iSq(i int) int {
	return i * i
}

func TestLists(t *testing.T) {
	tt := map[string]struct {
		expected []int
		list.Wrapper[int]
	}{
		"Generator": {
			expected: []int{0, 1, 4, 9, 16},
			Wrapper:  list.GeneratorFactory(iSq)(5).Wrap(),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			w := tc.Wrapper
			assert.Equal(t, len(tc.expected), w.Len())
			for i, e := range tc.expected {
				assert.Equal(t, e, tc.List.AtIdx(i))
			}
		})
	}
}

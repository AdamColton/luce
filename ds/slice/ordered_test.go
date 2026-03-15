package slice_test

import (
	"cmp"
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestOrdered(t *testing.T) {
	s := slice.New([]int{27, 13, 77, 30, 58, 15, 22, 53, 80, 8, 9, 25, 74, 40, 94, 34, 59, 69, 44, 71})
	o := slice.NewOrdered(s)

	expected := slice.New([]int{8, 9, 13, 15, 22, 25, 27, 30, 34, 40, 44, 53, 58, 59, 69, 71, 74, 77, 80, 94})
	assert.Equal(t, expected, o.Slice)

	validateO := func() {
		counter := 0
		o.Compare = slice.NewCompare(func(a, b int) int {
			counter++
			return cmp.Compare(a, b)
		})
		idx, found := o.Find(25)
		assert.Equal(t, 5, idx)
		assert.True(t, found)
		assert.Less(t, counter, 6)

		counter = 0
		idx, found = o.Find(70)
		assert.Equal(t, 15, idx)
		assert.False(t, found)
		assert.Less(t, counter, 6)

		assert.True(t, o.Contains(80))
		assert.False(t, o.Contains(10))
	}
	validateO()

	o = s.Ordered(o.Compare)
	validateO()
}

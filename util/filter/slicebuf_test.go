package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestFilterSliceBuf(t *testing.T) {
	gt2 := filter.GT(2)
	vals := []int{1, 5, 2, 6}

	assert.Equal(t, []int{5, 6}, []int(gt2.Slice(vals, nil)))

	// The buffer is used when it has the capacity.
	buf := make([]int, 0, 4)
	got := gt2.Slice(vals, buf)
	assert.Equal(t, []int{5, 6}, []int(got))
	assert.Equal(t, buf[:1][0], got[0])
	assert.Same(t, &buf[:1][0], &got[0])

	// The buffer can be the input.
	in := []int{1, 5, 2, 6}
	assert.Equal(t, []int{5, 6}, []int(gt2.Slice(in, in)))

	assert.Nil(t, filter.GT(9).Slice(vals, nil))
}

package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestByIdx(t *testing.T) {
	src := slice.New([]int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3})
	idxs := list.NewGenerator(src.Len()/2, func(i int) int {
		return i * 2
	})
	got := list.NewByIdx(src, idxs).Wrap().Slice(nil)
	expected := []int{3, 4, 5, 2, 5}
	assert.Equal(t, expected, got)
}

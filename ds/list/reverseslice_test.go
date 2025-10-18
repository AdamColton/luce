package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/stretchr/testify/assert"
)

func TestReverseSlice(t *testing.T) {
	r := list.ReverseSlice([]int{1, 2, 3})
	assert.Equal(t, []int{3, 2, 1}, []int(r.Slice(nil)))
}

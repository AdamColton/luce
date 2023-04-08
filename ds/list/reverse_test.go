package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/stretchr/testify/assert"
)

func TestNewReverse(t *testing.T) {
	r := list.NewReverse(list.Slice([]int{1, 2, 3}))
	assert.Equal(t, 3, r.Len())
	assert.Equal(t, 3, r.AtIdx(0))
	assert.Equal(t, 1, r.AtIdx(2))
}

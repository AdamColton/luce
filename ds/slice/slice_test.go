package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	cp := slice.Clone(data)
	assert.Equal(t, data, cp)
	data[0] = 0
	assert.Equal(t, 3, cp[0])
}

func TestSwap(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	slice.Swap(data, 0, 1)
	assert.Equal(t, 1, data[0])
	assert.Equal(t, 3, data[1])

}

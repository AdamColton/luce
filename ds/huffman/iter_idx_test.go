package huffman_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestIterIdx(t *testing.T) {
	tree := huffman.New([]huffman.Frequency[string]{{Val: "a", Count: 5}, {Val: "b", Count: 2}, {Val: "c", Count: 1}, {Val: "d", Count: 1}})
	data := slice.Slice[string]{"a", "b", "c"}
	bits := huffman.Encode[string](list.Wrap[string](data), huffman.NewLookup(tree))

	it := tree.Iter(bits)
	v, done := it.Cur()
	assert.Equal(t, "a", v)
	assert.False(t, done)
	assert.Equal(t, 0, it.Idx())

	v, done = it.Next()
	assert.Equal(t, "b", v)
	assert.False(t, done)
	assert.Equal(t, 1, it.Idx())

	v, done = it.Next()
	assert.Equal(t, "c", v)
	assert.False(t, done)
	assert.Equal(t, 2, it.Idx())

	// The last value has been returned, so Cur reports done.
	_, done = it.Cur()
	assert.True(t, done)
	_, done = it.Next()
	assert.True(t, done)
}

package huffman_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/huffman"
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestNewEmptyPanics(t *testing.T) {
	assert.Panics(t, func() { huffman.New[string](nil) })
	assert.Panics(t, func() { huffman.MapNew(map[string]int{}) })
}

func TestEncodeUnknownPanics(t *testing.T) {
	tree := huffman.New([]huffman.Frequency[string]{{Val: "a", Count: 2}, {Val: "b", Count: 1}})
	lookup := huffman.NewLookup(tree)
	assert.Nil(t, lookup.Get("z"))
	data := slice.Slice[string]{"a", "z"}
	assert.Panics(t, func() { huffman.Encode[string](list.Wrap[string](data), lookup) })
}

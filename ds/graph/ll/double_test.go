package ll_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/ds/graph/ll"
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestDouble(t *testing.T) {
	expected := slice.Slice[graph.KV[string, int]]{
		{"3", 3},
		{"1", 1},
		{"4", 4},
		{"1", 1},
		{"5", 5},
		{"9", 9},
	}
	n := ll.NewDoubleLoop(expected[0].K, expected[0].V)
	assert.Equal(t, n.Next, n)
	assert.Equal(t, n.Prev, n)
	s := n
	for _, kv := range expected[1:] {
		nxt := n.InsertAfter(kv.K, kv.V)
		assert.Equal(t, n.Next, nxt)
		assert.Equal(t, n, nxt.Prev)
		n = nxt
	}

	got := slice.New(slice.FromIter(s.Iter(true), nil))
	assert.Equal(t, expected, got)

	got = slice.New(slice.FromIter(n.Iter(false), nil))
	expected = slice.FromIter(list.NewReverse(expected).Iter(), nil)
	assert.Equal(t, expected, got)
}

package ll_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/ds/graph/ll"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestSinglePtr(t *testing.T) {
	expected := []graph.KV[string, int]{
		{"3", 3},
		{"1", 1},
		{"4", 4},
		{"1", 1},
		{"5", 5},
		{"9", 9},
	}
	p := graph.RawPointer[ll.Single[string, int]]{}
	n := ll.NewSingleLoop(p, expected[0].K, expected[0].V)
	s := n
	assert.Equal(t, n.Next.Get(), n)
	for _, kv := range expected[1:] {
		nxt := n.InsertAfter(kv.K, kv.V)
		assert.Equal(t, n.Next.Get(), nxt)
		n = nxt
	}

	got := slice.IterSlice(s.Iter(), nil)
	assert.Equal(t, expected, got)
}

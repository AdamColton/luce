package rbtree

import (
	"testing"

	"github.com/adamcolton/luce/ds/graph"
	"github.com/stretchr/testify/assert"
)

func TestNodeColor(t *testing.T) {
	n := &node[string, int]{
		color: red,
		size:  1,
		chld: [2]graph.Ptr[*node[string, int]]{
			graph.RawPointer[node[string, int]]{},
			graph.RawPointer[node[string, int]]{},
		},
		prt: graph.RawPointer[node[string, int]]{},
	}

	assert.Equal(t, red, n.clr())
	assert.Equal(t, black, n.getChild(0).clr())

}

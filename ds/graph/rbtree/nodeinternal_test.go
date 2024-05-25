package rbtree

import (
	"testing"

	"github.com/adamcolton/luce/ds/graph"
	"github.com/stretchr/testify/assert"
)

func TestNodeColor(t *testing.T) {
	n := &Node[string, int]{
		color: red,
		size:  1,
		chld: [2]graph.Ptr[*Node[string, int]]{
			graph.RawPointer[Node[string, int]]{},
			graph.RawPointer[Node[string, int]]{},
		},
		prt: graph.RawPointer[Node[string, int]]{},
	}

	assert.Equal(t, red, n.clr())
	assert.Equal(t, black, n.getChild(0).clr())

}

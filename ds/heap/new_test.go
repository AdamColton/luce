package heap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/heap"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	h := heap.New(func(a, b string) bool { return len(a) < len(b) })
	for _, s := range []string{"ccc", "a", "bb"} {
		h.Push(s)
	}
	assert.Equal(t, "a", h.Pop())
	assert.Equal(t, "bb", h.Pop())
	assert.Equal(t, "ccc", h.Pop())
}

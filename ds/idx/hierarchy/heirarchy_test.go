package hierarchy_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/hierarchy"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/stretchr/testify/assert"
)

func TestHierarchy(t *testing.T) {
	h := hierarchy.New[uint32, string](100)
	id, found := h.Get([]string{"this", "is", "a", "test"}, true)
	assert.False(t, found)
	k, found := h.B(id)
	assert.True(t, found)
	assert.Equal(t, "test", k.Name)

	h.Get([]string{"this", "test"}, true)
	expected := lset.New("is", "test")
	id, found = h.Get([]string{"this"}, false)
	assert.True(t, found)
	assert.Equal(t, expected, h.Children[id])
}

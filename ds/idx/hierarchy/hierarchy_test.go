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

func TestGet(t *testing.T) {
	h := hierarchy.New[uint8, string](0)

	// The root is 0.
	id, found := h.Get(nil, false)
	assert.Equal(t, uint8(0), id)
	assert.True(t, found)

	// Nothing is created without create.
	id, found = h.Get([]string{"a", "b"}, false)
	assert.Equal(t, uint8(0), id)
	assert.False(t, found)
	assert.Equal(t, uint8(1), h.MaxID)

	// IDs are assigned in order and found is false when the name is created.
	id, found = h.Get([]string{"a", "b"}, true)
	assert.Equal(t, uint8(2), id)
	assert.False(t, found)
	a, found := h.Get([]string{"a"}, false)
	assert.Equal(t, uint8(1), a)
	assert.True(t, found)

	// Getting it again finds it.
	id, found = h.Get([]string{"a", "b"}, true)
	assert.Equal(t, uint8(2), id)
	assert.True(t, found)

	// A missing name in the middle of a path.
	id, found = h.Get([]string{"a", "x", "y"}, false)
	assert.Equal(t, uint8(0), id)
	assert.False(t, found)
	assert.Equal(t, uint8(3), h.MaxID)
}

func TestKey(t *testing.T) {
	h := hierarchy.New[uint8, string](0)
	id, found := h.Key(0, "a", true)
	assert.Equal(t, uint8(1), id)
	assert.False(t, found)

	// The same name under a different parent has a different ID.
	other, _ := h.Key(0, "b", true)
	nested, found := h.Key(other, "a", true)
	assert.False(t, found)
	assert.NotEqual(t, id, nested)

	id2, found := h.Key(0, "a", false)
	assert.Equal(t, id, id2)
	assert.True(t, found)

	_, found = h.Key(0, "missing", false)
	assert.False(t, found)

	// The parent must be in the Hierarchy to create a name under it.
	assert.Panics(t, func() { h.Key(99, "z", true) })
}

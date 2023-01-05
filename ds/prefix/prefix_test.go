package prefix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	p := New()

	n, insert := p.Upsert("")
	assert.Nil(t, n)
	assert.False(t, insert)

	n, insert = p.Upsert("abc")
	assert.True(t, n.IsWord())
	assert.True(t, insert)
	assert.Equal(t, "abc", n.Gram())
	assert.Equal(t, 0, n.ChildrenCount())
	_, insert = p.Upsert("abc")
	assert.False(t, insert)

	n = p.root.children.GetVal('a')
	assert.False(t, n.IsWord())
	assert.Equal(t, "a", n.Gram())
	assert.Equal(t, 1, n.ChildrenCount())
	assert.Equal(t, []rune{'b'}, n.Children())

	n = n.Child('b')
	assert.False(t, n.IsWord())
	assert.Equal(t, "ab", n.Gram())
	assert.Equal(t, 1, n.ChildrenCount())
	assert.Equal(t, []rune{'c'}, n.Children())
}

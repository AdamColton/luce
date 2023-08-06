package navigator_test

import (
	"testing"

	"github.com/adamcolton/luce/util/navigator"
	"github.com/stretchr/testify/assert"
)

type node struct {
	children map[rune]*node
}

func (n *node) Next(key rune, create bool, ctx navigator.VoidContext) (next *node, ok bool) {
	next, ok = n.children[key]
	if !ok && create {
		next = &node{
			children: map[rune]*node{},
		}
		n.children[key] = next
		ok = true
	}
	return
}

func TestNavigator(t *testing.T) {
	root := &node{
		children: map[rune]*node{},
	}
	n := navigator.New(root, []rune("abc"))

	assert.Equal(t, 'a', n.IdxKey())

	got, ok := n.Seek(false, navigator.Void)
	assert.False(t, ok)
	assert.Nil(t, got)

	n.Cur = root
	got, ok = n.Seek(true, navigator.Void)
	assert.True(t, ok)
	expected := root.children['a'].children['b'].children['c']
	assert.Equal(t, expected, got)

	n.Cur = root
	n.Idx = 0
	got, ok = n.Trace(true).Seek(false, navigator.Void)
	assert.True(t, ok)
	assert.Equal(t, expected, got)

	got, ok = n.Pop()
	assert.True(t, ok)
	assert.Equal(t, expected, got)
	expected = root.children['a'].children['b']
	assert.Equal(t, expected, n.Cur)
}

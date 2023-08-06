package navigator_test

import (
	"testing"

	"github.com/adamcolton/luce/util/navigator"
	"github.com/stretchr/testify/assert"
)

func TestSeekTrace(t *testing.T) {
	root := &node{children: map[rune]*node{}}
	n := navigator.New(root, []rune("ab")).Trace(true)
	got, ok := n.Seek(true, navigator.Void)
	assert.True(t, ok)
	assert.Equal(t, got, n.Cur)
	assert.Equal(t, 2, n.Idx)

	// Nodes holds every node visited except the last, which is Cur.
	a := root.children['a']
	assert.Len(t, n.Nodes, 2)
	assert.Equal(t, root, n.Nodes[0])
	assert.Equal(t, a, n.Nodes[1])

	// Pop undoes the last step.
	popped, ok := n.Pop()
	assert.True(t, ok)
	assert.Equal(t, got, popped)
	assert.Equal(t, a, n.Cur)
	assert.Equal(t, 1, n.Idx)
	assert.Equal(t, 'b', n.IdxKey())

	// Seek continues from there.
	again, ok := n.Seek(false, navigator.Void)
	assert.True(t, ok)
	assert.Equal(t, got, again)
}

func TestSeekFailure(t *testing.T) {
	root := &node{children: map[rune]*node{}}
	root.children['a'] = &node{children: map[rune]*node{}}
	n := navigator.New(root, []rune("abc")).Trace(true)

	got, ok := n.Seek(false, navigator.Void)
	assert.False(t, ok)
	assert.Nil(t, got)

	// Cur is what Next returned and Idx is the key that failed.
	assert.Nil(t, n.Cur)
	assert.Equal(t, 1, n.Idx)
	assert.Equal(t, 'b', n.IdxKey())
	assert.Equal(t, []*node{root, root.children['a']}, []*node(n.Nodes))
}

func TestSeekNoKeys(t *testing.T) {
	root := &node{children: map[rune]*node{}}
	n := navigator.New(root, nil)
	got, ok := n.Seek(false, navigator.Void)
	assert.True(t, ok)
	assert.Equal(t, root, got)
}

func TestPopEmpty(t *testing.T) {
	root := &node{children: map[rune]*node{}}
	n := navigator.New(root, []rune("a"))
	got, ok := n.Pop()
	assert.False(t, ok)
	assert.Nil(t, got)
	assert.Equal(t, root, n.Cur)
	assert.Equal(t, 0, n.Idx)
}

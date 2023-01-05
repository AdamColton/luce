package prefix

import (
	"github.com/adamcolton/luce/ds/lmap"
)

// Node in a prefix tree.
type Node interface {
	// Child looks up the child node by rune
	Child(rune) Node
	// Children returns all the child runes
	Children() []rune
	// ChildrenCount returns the number of children
	ChildrenCount() int
	// IsWord returns true if this Node was inserted as word into the tree
	IsWord() bool
	// Gram returns the string this node represents
	Gram() string
}

type node struct {
	isWord        bool
	r             rune
	parent        *node
	children      lmap.Wrapper[rune, *node]
	childrenCount int
}

func newNode() *node {
	return &node{
		children: lmap.New[rune, *node](nil),
	}
}

func (n *node) IsWord() bool {
	return n.isWord
}

func (n *node) ChildrenCount() int {
	return n.childrenCount
}

func (n *node) Gram() string {
	return string(n.recursiveGram(0))
}

func (n *node) recursiveGram(ln int) []rune {
	if n.parent == nil {
		return make([]rune, 0, ln)
	}
	return append(n.parent.recursiveGram(ln+1), n.r)
}

func (n *node) Child(r rune) Node {
	return n.children.GetVal(r)
}

func (n *node) Children() []rune {
	// TODO: use buf
	return n.children.Keys(nil)
}

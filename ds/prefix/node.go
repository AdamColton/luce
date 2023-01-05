package prefix

import "github.com/adamcolton/luce/ds/slice"

type node struct {
	isWord        bool
	r             rune
	parent        *node
	children      map[rune]*node
	childrenCount int
}

type Node interface {
	Child(rune) Node
	Children() []rune
	ChildrenCount() int
	IsWord() bool
	Gram() string
	Suggest(max int) []Suggestion
}

func newNode() *node {
	return &node{
		children: make(map[rune]*node),
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
	return n.children[r]
}

func (n *node) Children() []rune {
	return slice.Keys(n.children)
}

package rbtree

import (
	"fmt"

	"github.com/adamcolton/luce/math/cmpr"
)

func (c color) String() string {
	if c == -1 {
		return "Red"
	} else if c == 1 {
		return "Black"
	}
	return "???"
}

func (n *Node[Key, Val]) print() {
	if n == nil {
		return
	}
	fmt.Print("(")
	n.getChild(0).print()
	fmt.Print(n.color.String(), " ", n.K, " ", n.size)
	n.getChild(1).print()
	fmt.Print(")")
}

func (n *Node[Key, Val]) Validate() (blackSum int, ok bool) {
	if n == nil {
		return 1, true
	}
	c0 := n.getChild(0)
	c1 := n.getChild(1)
	b0, ok0 := c0.Validate()
	b1, ok1 := c1.Validate()
	sizeOk := n.size == c0.getSize()+c1.getSize()+1
	ok = ok0 && ok1 && b0 == b1 && sizeOk
	blackSum = cmpr.Max(b0, b1)
	if n.color == red {
		ok = ok && c0.clr() == black && c1.clr() == black
	} else {
		blackSum++
	}
	return
}

func (t *Tree[Key, Val]) Print() {
	r := stripBool(t.root.Get())
	if r != nil {
		r.print()
	} else {
		fmt.Print("()")
	}
	fmt.Print("\n")
}

func (t *Tree[Key, Val]) Validate() bool {
	r := stripBool(t.root.Get())
	if r == nil {
		return true
	}
	_, ok := r.Validate()
	return ok
}

package prefix

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/entity"
)

type NodeBale struct {
	IsWord        bool
	R             rune
	ChildrenCount int
	Children      map[rune]*NodeBale
}

var (
	// set in init to avoid initilization cycle error
	baleChildren   lmap.TransformFunc[rune, *node, rune, *NodeBale]
	unbaleChildren lmap.TransformFunc[rune, *NodeBale, rune, *node]
)

func (n *node) bale() *NodeBale {
	return &NodeBale{
		IsWord:        n.isWord,
		R:             n.r,
		ChildrenCount: n.childrenCount,
		Children:      baleChildren.Transform(n.children, nil).Map(),
	}
}

func (bale *NodeBale) unbale() *node {
	n := &node{}
	bale.unbaleTo(n)
	n.children.Each(func(r rune, child *node, done *bool) {
		child.parent = n
	})
	return n
}

func (bale *NodeBale) unbaleTo(n *node) {
	n.isWord = bale.IsWord
	n.r = bale.R
	n.childrenCount = bale.ChildrenCount
	n.children = unbaleChildren.Map(bale.Children)
}

type PrefixBale struct {
	Root *NodeBale
}

func (p *Prefix) saveIf() {
	if p.save {
		p.Save()
	}
}

func (p *Prefix) Bale() *PrefixBale {
	gp := &PrefixBale{
		Root: p.root.bale(),
	}
	return gp
}

func (bale *PrefixBale) Unbale() *Prefix {
	p := &Prefix{}
	bale.UnbaleTo(p)
	return p
}

func (bale *PrefixBale) EntRefs() []entity.Key {
	return nil
}

func (bale *PrefixBale) UnbaleTo(p *Prefix) {
	p.root = bale.Root.unbale()
	p.starts = make(map[rune]slice.Slice[*node])
	p.root.populateStarts(p)
}

func (n *node) populateStarts(p *Prefix) {
	n.children.Each(func(r rune, child *node, done *bool) {
		p.starts[r] = append(p.starts[r], child)
		child.populateStarts(p)
	})
}

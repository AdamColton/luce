package prefix

import (
	"sort"
)

type Suggestion struct {
	Word      string
	Terminals []int
}

func (n *node) Suggest(max int) []Suggestion {
	type childCount struct {
		r rune
		c int
	}
	ccs := make([]childCount, 0, len(n.children))
	for _, c := range n.children {
		ccs = append(ccs, childCount{
			r: c.r,
			c: c.childrenCount,
		})
	}
	sort.Slice(ccs, func(i, j int) bool {
		return ccs[i].c > ccs[j].c
	})
	if len(ccs) > max {
		ccs = ccs[:max]
	}

	out := make([]Suggestion, 0, len(ccs))
	for _, cc := range ccs {
		c := n.children[cc.r]
		rs, terminals := c.suggestion(0)
		out = append(out, Suggestion{
			Word:      string(rs),
			Terminals: terminals,
		})
	}

	return out
}

func (n *node) suggestion(ln int) (rs []rune, terminals []int) {
	var best *node
	for _, c := range n.children {
		if best == nil || c.childrenCount > best.childrenCount {
			best = c
		}
	}

	if best == nil {
		rs = make([]rune, ln+1)
		terminals = []int{ln}
	} else {
		rs, terminals = best.suggestion(ln + 1)
	}
	rs[ln] = n.r
	if n.isWord {
		terminals = append(terminals, ln)
	}
	return rs, terminals
}

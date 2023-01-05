package prefix

import (
	"sort"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
)

type Suggestion struct {
	Word      string
	Terminals slice.Slice[int]
}

func (s Suggestion) Words(prefix string) list.Wrapper[string] {
	return list.NewTransformer(s.Terminals, func(term int) string {
		return prefix + s.Word[:term+1]
	})
}

func (n *node) Suggest(max int) []Suggestion {
	type childCount struct {
		r rune
		c int
	}
	ccs := make([]childCount, 0, n.children.Len())
	n.children.Each(func(key rune, c *node, done *bool) {
		ccs = append(ccs, childCount{
			r: c.r,
			c: c.childrenCount,
		})
	})
	sort.Slice(ccs, func(i, j int) bool {
		return ccs[i].c > ccs[j].c
	})
	if len(ccs) > max {
		ccs = ccs[:max]
	}

	out := make([]Suggestion, 0, len(ccs))
	for _, cc := range ccs {
		c := n.children.GetVal(cc.r)
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
	n.children.Each(func(key rune, child *node, done *bool) {
		if best == nil || child.childrenCount > best.childrenCount {
			best = child
		}
	})

	if best == nil {
		rs = make([]rune, ln+1)
	} else {
		rs, terminals = best.suggestion(ln + 1)
	}
	rs[ln] = n.r
	if n.isWord {
		terminals = append(terminals, ln)
	}
	return rs, terminals
}

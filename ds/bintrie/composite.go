package bintrie

func Or[U Uint](a, b Trie[U]) Trie[U] {
	return or(a.(*node[U]), b.(*node[U]))
}

func or[U Uint](x, y *node[U]) *node[U] {
	t := x.terminal || y.terminal
	b0 := x.branches[0] != nil || y.branches[0] != nil
	b1 := x.branches[1] != nil || y.branches[1] != nil
	out := &node[U]{
		terminal: t,
	}
	if b0 {
		if x.branches[0] == nil {
			out.branches[0] = y.branches[0].copy()
		} else if y.branches[0] == nil {
			out.branches[0] = x.branches[0].copy()
		} else {
			out.branches[0] = or(x.branches[0], y.branches[0])
		}
	}
	if b1 {
		if x.branches[1] == nil {
			out.branches[1] = y.branches[1].copy()
		} else if y.branches[1] == nil {
			out.branches[1] = x.branches[1].copy()
		} else {
			out.branches[1] = or(x.branches[1], y.branches[1])
		}
	}
	out.updateSize()
	return out
}

func And[U Uint](a, b Trie[U]) Trie[U] {
	return and(a.(*node[U]), b.(*node[U]))
}

func and[U Uint](x, y *node[U]) *node[U] {
	t := x.terminal && y.terminal
	and0 := x.branches[0] != nil && y.branches[0] != nil
	and1 := x.branches[1] != nil && y.branches[1] != nil
	if !t && !and0 && !and1 {
		return nil
	}
	n := &node[U]{
		terminal: t,
	}
	if and0 {
		n.branches[0] = and(x.branches[0], y.branches[0])
	}
	if and1 {
		n.branches[1] = and(x.branches[1], y.branches[1])
	}
	n.updateSize()
	return n.nilCheck()
}

// Nand returns all values in a but not in b.
func Nand[U Uint](a, b Trie[U]) Trie[U] {
	return nand(a.(*node[U]), b.(*node[U]))
}

func nand[U Uint](x, y *node[U]) *node[U] {
	t := x.terminal && !y.terminal
	x0 := x.branches[0] != nil
	x1 := x.branches[1] != nil
	if !t && !x0 && !x1 {
		return nil
	}
	n := &node[U]{}
	if x0 {
		if y.branches[0] != nil {
			n.branches[0] = nand(x.branches[0], y.branches[0])
		} else {
			n.branches[0] = x.branches[0].copy()
		}
	}
	if x1 {
		if y.branches[1] != nil {
			n.branches[1] = nand(x.branches[1], y.branches[1])
		} else {
			n.branches[1] = x.branches[1].copy()
		}
	}
	n.updateSize()
	return n.nilCheck()
}

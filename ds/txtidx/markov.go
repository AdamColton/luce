package txtidx

type Markov struct {
	root  *MarkovNode
	heads map[rune][]*MarkovNode
}

func NewMarkov() *Markov {
	return &Markov{
		root: &MarkovNode{
			Next: map[rune]*MarkovNode{},
		},
		heads: map[rune][]*MarkovNode{},
	}
}

func (m *Markov) Find(str string) *Word {
	n := m.root.Find([]rune(str))
	if n == nil {
		return nil
	}
	return n.Word
}

func (m *Markov) FindAll(str string) []*Word {
	rs := []rune(str)
	if len(rs) == 0 {
		return nil
	}
	var out []*Word
	for _, r := range m.heads[rs[0]] {
		n := r.Find(rs[1:])
		if n != nil && n.Word != nil {
			out = append(out, n.Word)
		}
	}
	return out
}

func (m *Markov) Upsert(str string) *Word {
	n := m.root.Upsert([]rune(str), m)
	if n.Word == nil {
		n.Word = &Word{
			WordID:    WordID(MaxUint32),
			Documents: newDocSet(),
		}
	}
	return n.Word
}

type MarkovNode struct {
	Next map[rune]*MarkovNode
	*Word
}

func (n *MarkovNode) Find(rs []rune) *MarkovNode {
	if n == nil || len(rs) == 0 {
		return n
	}
	return n.Next[rs[0]].Find(rs[1:])
}

func (n *MarkovNode) Upsert(rs []rune, root *Markov) *MarkovNode {
	if len(rs) == 0 {
		return n
	}
	r, rs := rs[0], rs[1:]
	nn := n.Next[r]
	if nn == nil {
		nn = &MarkovNode{
			Next: map[rune]*MarkovNode{},
		}
		n.Next[r] = nn
		root.heads[r] = append(root.heads[r], nn)
	}
	return nn.Upsert(rs, root)
}

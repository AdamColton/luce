package txtidx

type Markov struct {
	heads map[rune][]*MarkovNode
	nodes map[mkey]*MarkovNode
	maxID uint32
}

type mkey struct {
	nodeID uint32
	r      rune
}

func NewMarkov() *Markov {
	return &Markov{
		heads: map[rune][]*MarkovNode{},
		nodes: map[mkey]*MarkovNode{},
	}
}

func (m *Markov) Find(str string) *Word {
	s := newSeeker(str)
	s.find(m.nodes)
	if s.n == nil {
		return nil
	}
	return s.n.Word
}

func (m *Markov) FindAll(str string) []*Word {
	rs := []rune(str)
	if len(rs) == 0 {
		return nil
	}

	s := newSeeker(str)

	var out []*Word
	for _, n := range m.heads[s.rs[0]] {
		s.n = n
		s.idx = 1
		s.find(m.nodes)
		if s.n != nil {
			out = m.AppendAll(s.n, out)
		}
	}
	return out
}

type seeker struct {
	rs  []rune
	idx int
	k   mkey
	n   *MarkovNode
}

func newSeeker(str string) *seeker {
	s := &seeker{
		rs: []rune(str),
	}
	return s
}

func (s *seeker) moveNext() (notDone bool) {
	notDone = s.idx < len(s.rs)
	if notDone {
		s.k.r = s.rs[s.idx]
		if s.n != nil {
			s.k.nodeID = s.n.id
		}
		s.idx++
	}
	return notDone
}

func (s *seeker) find(m map[mkey]*MarkovNode) {
	for s.moveNext() {
		s.n = m[s.k]
		if s.n == nil {
			return
		}
	}
}

func (m *Markov) Upsert(str string) *Word {
	if len(str) == 0 {
		panic("don't do that")
	}
	s := newSeeker(str)
	for s.moveNext() {
		next := m.nodes[s.k]
		if next == nil {
			if s.n != nil {
				s.n.children = append(s.n.children, s.k.r)
			}
			m.maxID++
			next = &MarkovNode{
				id: m.maxID,
			}
			m.nodes[s.k] = next
			m.heads[s.k.r] = append(m.heads[s.k.r], next)
		}
		s.n = next
	}
	if s.n.Word == nil {
		s.n.Word = &Word{
			WordIDX:   WordIDX(MaxUint32),
			Documents: newDocSet(),
		}
	}
	return s.n.Word
}

func (m *Markov) AppendAll(n *MarkovNode, buf []*Word) []*Word {
	if n.Word != nil {
		buf = append(buf, n.Word)
	}
	k := mkey{
		nodeID: n.id,
	}
	for _, r := range n.children {
		k.r = r
		buf = m.AppendAll(m.nodes[k], buf)
	}
	return buf
}

type MarkovNode struct {
	id       uint32
	children []rune
	*Word
}

package txtidx

import "sort"

type markov struct {
	nodes map[mkey]*markovNode
	maxID uint32
}

type mkey struct {
	nodeID uint32
	r      rune
}

func newMarkov() *markov {
	return &markov{
		nodes: map[mkey]*markovNode{},
	}
}

func (m *markov) find(str string) (*word, *markovNode) {
	s := newSeeker(str)
	s.find(m.nodes)
	if s.n == nil {
		return nil, nil
	}
	return s.n.word, s.n
}

func (m *markov) findAll(str string) words {
	s := newSeeker(str)
	if s == nil {
		return nil
	}
	s.find(m.nodes)
	if s.n == nil {
		return nil
	}

	return m.appendAll(s.n, nil)
}

type seeker struct {
	rs  []rune
	idx int
	k   mkey
	n   *markovNode
}

func newSeeker(str string) *seeker {
	rs := []rune(str)
	if len(rs) == 0 {
		return nil
	}
	s := &seeker{
		rs: rs,
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

func (s *seeker) find(m map[mkey]*markovNode) {
	for s.moveNext() {
		s.n = m[s.k]
		if s.n == nil {
			return
		}
	}
}

func (m *markov) upsert(str string) *word {
	if len(str) == 0 {
		panic("don't do that")
	}
	w, _ := m.upsertRecursive(newSeeker(str))
	return w
}

func (m *markov) upsertRecursive(s *seeker) (*word, bool) {
	if !s.moveNext() {
		inc := s.n.word == nil
		if inc {
			s.n.word = &word{
				wordIDX:   wordIDX(MaxUint32),
				Documents: newDocSet(),
			}
		}
		return s.n.word, inc
	}
	n := s.n
	r := s.k.r
	next := m.nodes[s.k]
	if next == nil {
		if n != nil {
			n.children = append(n.children, childCount{
				r: r,
			})
		}
		m.maxID++
		next = &markovNode{
			id: m.maxID,
		}
		m.nodes[s.k] = next
	}
	s.n = next
	w, inc := m.upsertRecursive(s)
	if inc && n != nil {
		i := findChild(r, n.children)
		n.children[i].count++
	}
	return w, inc
}

func (m *markov) appendAll(n *markovNode, buf []*word) []*word {
	if n.word != nil {
		buf = append(buf, n.word)
	}
	k := mkey{
		nodeID: n.id,
	}
	for _, c := range n.children {
		k.r = c.r
		buf = m.appendAll(m.nodes[k], buf)
	}
	return buf
}

func (m *markov) deleteWord(w string) {
	rs := []rune(w)
	m.recursiveDeleteWord(m.nodes[mkey{
		r: rs[0],
	}], rs[1:])
}

func (m *markov) recursiveDeleteWord(n *markovNode, rs []rune) (bool, bool) {
	if n == nil {
		return false, false
	}
	if len(rs) == 0 {
		n.word = nil
		return len(n.children) == 0, true
	}
	k := mkey{
		r:      rs[0],
		nodeID: n.id,
	}
	shouldDelete, shouldDec := m.recursiveDeleteWord(m.nodes[k], rs[1:])
	if shouldDelete {
		delete(m.nodes, k)
		ln := len(n.children)
		if ln == 1 {
			return true, true
		}
		ln--
		i := findChild(rs[0], n.children[:ln])
		if i >= 0 {
			n.children[i] = n.children[ln]
		}
		n.children = n.children[:ln]
	} else if shouldDec {
		i := findChild(rs[0], n.children)
		n.children[i].count--
	}
	return false, shouldDec
}

type childCount struct {
	r     rune
	count uint16
}

func findChild(r rune, cs []childCount) int {
	for i, c := range cs {
		if c.r == r {
			return i
		}
	}
	return -1
}

type markovNode struct {
	id       uint32
	children []childCount
	word     *word
}

func (mn *markovNode) sort() {
	sort.Slice(mn.children, func(i, j int) bool {
		return mn.children[i].count > mn.children[j].count
	})
}

func (m *markov) suggest(word string, max int) []Suggestion {
	_, n := m.find(word)
	n.sort()
	if ln := len(n.children); ln < max || max < 0 {
		max = ln
	}
	out := make([]Suggestion, max)
	for i := range n.children[:max] {
		out[i] = m.expand(n, i)
	}
	return out
}

func (m *markov) expand(n *markovNode, cIdx int) Suggestion {
	k := mkey{
		nodeID: n.id,
		r:      n.children[cIdx].r,
	}
	out, terminals := m.expandRecursive(m.nodes[k], 1)
	out[0] = k.r
	return Suggestion{
		Word:      string(out),
		Terminals: terminals,
	}
}

func (m *markov) expandRecursive(n *markovNode, d int) ([]rune, []int) {
	if len(n.children) == 0 {
		var terminals []int
		if n.word != nil {
			terminals = []int{d}
		}
		return make([]rune, d), terminals
	}
	n.sort()
	k := mkey{
		nodeID: n.id,
		r:      n.children[0].r,
	}
	next := m.nodes[k]
	out, terminals := m.expandRecursive(next, d+1)
	out[d] = n.children[0].r
	if n.word != nil {
		terminals = append(terminals, d)
	}
	return out, terminals
}

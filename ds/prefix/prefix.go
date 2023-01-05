package prefix

import "github.com/adamcolton/luce/ds/slice"

type Prefix struct {
	root   *node
	starts map[rune][]*node
}

func New() *Prefix {
	return &Prefix{
		root:   newNode(),
		starts: make(map[rune][]*node),
	}
}

func (p *Prefix) seeker(str string) *seeker {
	rs := []rune(str)
	s := &seeker{
		runes: rs,
		p:     p,
		n:     p.root,
	}
	return s
}

func (p *Prefix) Upsert(word string) (Node, bool) {
	if len(word) == 0 {
		return nil, false
	}
	s := p.seeker(word)
	for done := s.moveNext(true); !done; done = s.moveNext(true) {
	}
	insert := false
	if !s.n.isWord {
		insert = true
		s.n.isWord = true
		for p := s.n.parent; p != nil; p = p.parent {
			p.childrenCount++
		}
	}
	return s.n, insert
}

func (p *Prefix) Find(gram string) Node {
	s := p.seeker(gram)
	for done := false; !done; done = s.moveNext(false) {
	}
	return s.n
}

func (p *Prefix) Contains(gram string) []Node {
	rs := []rune(gram)
	if len(rs) == 0 {
		return nil
	}
	s := &seeker{
		runes: rs[1:],
		p:     p,
	}
	var out []Node
	for _, n := range p.starts[rs[0]] {
		s.idx = 0
		s.n = n
		for done := false; !done; done = s.moveNext(false) {
		}
		out = slice.AppendNotNil[Node](out, s.n)
	}
	return out
}

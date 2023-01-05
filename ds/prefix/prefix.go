// Package prefix holds a prefix tree for strings. The term word is used to
// refer to values which are present in the tree. Gram is used to refer to any
// sequence that exists in the tree. For instance is the word 'test' is inserted
// then the tree contains the gram 'tes' but it will not be a word.
package prefix

import "github.com/adamcolton/luce/ds/slice"

// Prefix is the root of a prefix tree
type Prefix struct {
	root   *node
	starts map[rune]slice.Slice[*node]
	purge  map[rune]bool
}

// New Prefix tree.
func New() *Prefix {
	return &Prefix{
		root:   newNode(),
		starts: make(map[rune]slice.Slice[*node]),
		purge:  make(map[rune]bool),
	}
}

func (p *Prefix) seeker(str string) *seeker {
	s := &seeker{
		runes: []rune(str),
		p:     p,
		n:     p.root,
	}
	return s
}

// Upsert a word into the prefix tree returning the Node for that word and a
// bool indicating if an insert happened.
func (p *Prefix) Upsert(word string) (n Node, insert bool) {
	if len(word) == 0 {
		return nil, false
	}
	s := p.seeker(word)
	for done := s.moveNext(true); !done; done = s.moveNext(true) {
	}
	if !s.n.isWord {
		insert = true
		s.n.isWord = true
		for p := s.n.parent; p != nil; p = p.parent {
			p.childrenCount++
		}
	}
	n = s.n
	return
}

// Find a node by it's gram. If there are no prefixes starting with the gram,
// nil is returned.
func (p *Prefix) Find(gram string) Node {
	return p.find(gram).n
}

func (p *Prefix) find(gram string) *seeker {
	s := p.seeker(gram)
	for done := false; !done; done = s.moveNext(false) {
	}
	return s
}

// Containing returns all nodes in the tree containing the specified gram,
// even if they do not begin with the specified gram.
func (p *Prefix) Containing(gram string) Nodes {

	if len(p.purge) > 0 {
		p.Purge()
	}

	rs := []rune(gram)
	if len(rs) == 0 {
		return nil
	}
	s := &seeker{
		runes: rs[1:],
		p:     p,
	}

	var out slice.Slice[Node]
	starts := p.starts[rs[0]]
	for _, n := range starts {
		s.idx = 0
		s.n = n
		for done := false; !done; done = s.moveNext(false) {
		}
		out = out.AppendNotZero(s.n)
	}
	return Nodes(out)
}

func (p *Prefix) Remove(word string) {
	s := p.find(word)
	if s.n == nil || !s.n.isWord {
		return
	}
	s.n.isWord = false
	for n, done := s.movePrev(); !done; n, done = s.movePrev() {
		if len(n.children) > 0 || n.isWord {
			break
		}
		delete(n.parent.children, s.runes[s.idx])
		n.parent = nil
	}
	for _, r := range s.runes {
		p.purge[r] = true
	}

}

func (p *Prefix) Purge() {
	purged := make([]rune, 0, len(p.purge))
	for r := range p.purge {
		purged = append(purged, r)
		nodes := p.starts[r]
		changed := false
		for i := 0; i < len(nodes); i++ {
			n := nodes[i]
			if n.parent == nil {
				nodes = nodes.Remove(i)
				i--
				changed = true
			}
		}
		if changed {
			if len(nodes) == 0 {
				delete(p.starts, r)
			} else {
				p.starts[r] = nodes
			}
		}
	}
	for _, r := range purged {
		delete(p.purge, r)
	}
}

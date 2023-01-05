// Package prefix holds a prefix tree for strings. The term word is used to
// refer to values which are present in the tree. Gram is used to refer to any
// sequence that exists in the tree. For instance is the word 'test' is inserted
// then the tree contains the gram 'tes' but it will not be a word.
package prefix

// Prefix is the root of a prefix tree
type Prefix struct {
	root *node
}

// New Prefix tree.
func New() *Prefix {
	return &Prefix{
		root: newNode(),
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

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

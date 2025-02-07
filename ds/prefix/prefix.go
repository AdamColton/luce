// Package prefix holds a prefix tree for strings. The term word is used to
// refer to values which are present in the tree. Gram is used to refer to any
// sequence that exists in the tree. For instance is the word 'test' is inserted
// then the tree contains the gram 'tes' but it will not be a word.
package prefix

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/util/navigator"
)

// Prefix is the root of a prefix tree
type Prefix struct {
	key    entity.Key
	root   *node
	starts map[rune]slice.Slice[*node]
	purge  map[rune]bool
	save   bool
}

// New Prefix tree.
func New() *Prefix {
	return &Prefix{
		key:    entity.Rand(),
		root:   newNode(),
		starts: make(map[rune]slice.Slice[*node]),
		purge:  make(map[rune]bool),
	}
}

func (p *Prefix) seeker(str string) *seeker {
	s := &seeker{
		p: p,
		Navigator: &navigator.Navigator[rune, *node, *Prefix]{
			Cur:  p.root,
			Keys: []rune(str),
		},
	}
	return s
}

// Upsert a word into the prefix tree returning the Node for that word and a
// bool indicating if an insert happened.
func (p *Prefix) Upsert(word string) (n Node, insert bool) {
	if len(word) == 0 {
		return nil, false
	}
	nd, _ := p.seeker(word).Seek(true, p)
	if !nd.isWord {
		insert = true
		if nd.setIsWord(true) {
		}
		for prnt := nd.parent; prnt != nil; prnt = prnt.parent {
			prnt.setChildrenCount(prnt.childrenCount + 1)
		}
		p.saveIf()
	}
	n = nd
	return
}

// Find a node by it's gram. If there are no prefixes starting with the gram,
// nil is returned.
func (p *Prefix) Find(gram string) Node {
	n := p.find(gram, false).Cur
	if n == nil {
		return nil
	}
	return n
}

func (p *Prefix) find(gram string, trace bool) *seeker {
	s := p.seeker(gram)
	s.Trace(trace)
	s.Seek(false, p)
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
	//s := p.seeker(gram)
	s := &seeker{
		p: p,
		Navigator: &navigator.Navigator[rune, *node, *Prefix]{
			Keys: rs[1:],
		},
	}

	var out slice.Slice[Node]
	starts := p.starts[rs[0]]
	for _, n := range starts {
		s.Idx = 0
		s.Cur = n
		n, ok := s.Seek(false, p)
		if ok {
			out = append(out, n)
		}
	}
	return Nodes(out)
}

func (p *Prefix) Remove(word string) {
	s := p.find(word, true)
	if s.Cur == nil || !s.Cur.isWord {
		return
	}
	s.Cur.setIsWord(false)
	for n, ok := s.Pop(); ok; n, ok = s.Pop() {
		if n.children.Len() > 0 || n.isWord {
			break
		}
		n.parent.deleteChild(s.Keys[s.Idx])
		n.parent = nil
	}
	for _, r := range s.Keys {
		p.purge[r] = true
	}
	p.saveIf()
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

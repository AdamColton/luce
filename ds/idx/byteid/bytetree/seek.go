package bytetree

import "bytes"

type stackFrame struct {
	b byte
	n *node
}

type seekResult struct {
	*node
	idIdx int
	found bool
	stack []stackFrame
}

func (bt *byteIdxByteTree) seek(id []byte, stack bool) *seekResult {
	sr := &seekResult{
		node: bt.root,
	}
	for _, b := range id {
		child := sr.children[b]
		if child == nil {
			return sr
		}
		if stack {
			sr.stack = append(sr.stack, stackFrame{b: b, n: sr.node})
		}
		sr.node = child
		sr.idIdx++
		sr.found = sr.match(id)
		if sr.found {
			return sr
		}
	}
	return sr
}

func (sr *seekResult) match(id []byte) bool {
	return sr.idx != -1 && bytes.Equal(id[sr.idIdx:], sr.rest)
}

func (sr *seekResult) insert(id []byte, idx int) {
	for _, b := range id[sr.idIdx:] {
		if sr.idx != -1 && len(sr.rest) > 0 {
			sr.bump()
		}
		child := sr.children[b]
		if child == nil {
			sr.children[b] = newNode(idx, id[sr.idIdx+1:])
			sr.childCount++
			return
		}
		sr.node = child
		sr.idIdx++
	}
	if sr.idx != -1 {
		panic("something has gone terribly wrong")
	}
	sr.idx = idx
}

func newNode(idx int, rest []byte) *node {
	n := &node{
		rest: make([]byte, len(rest)),
		idx:  idx,
	}
	copy(n.rest, rest)
	return n
}

func (n *node) bump() {
	n.children[n.rest[0]] = newNode(n.idx, n.rest[1:])
	n.childCount++
	n.idx = -1
	n.rest = nil
}

func (sr *seekResult) del(id []byte) {
	sr.idx = -1
	sr.rest = nil
	// prune tree
	ln := len(sr.stack)
	for sr.idx == -1 && sr.childCount == 0 && ln > 0 {
		parent := sr.stack[ln-1].n
		sr.idIdx--
		parent.children[id[sr.idIdx]] = nil
		ln--
		sr.node = parent
	}
}

// rightThenUp searches for a non-nil child proceeding right across each
// set of children, then if none are found moving up to the parent and trying
// again.
func (sr *seekResult) rightThenUp(gt int) bool {
	ln := len(sr.stack)
	for {
		for i := gt + 1; i < 256; i++ {
			if sr.children[i] != nil {
				sr.stack = append(sr.stack[:ln], stackFrame{b: byte(i), n: sr.node})
				sr.node = sr.children[i]
				sr.idIdx = ln
				return true
			}
		}
		if ln == 0 {
			break
		}
		ln--
		gt = int(sr.stack[ln].b)
		sr.node = sr.stack[ln].n
	}
	return false
}

func (sr *seekResult) downAndLeft() {
	for sr.idx == -1 {
		for i, c := range sr.children {
			if c != nil {
				sr.stack = append(sr.stack, stackFrame{b: byte(i), n: sr.node})
				sr.node = c
				break
			}
		}
	}
}

func (sr *seekResult) value() []byte {
	out := make([]byte, len(sr.stack)+len(sr.rest))
	for i, sf := range sr.stack {
		out[i] = sf.b
	}
	copy(out[len(sr.stack):], sr.rest)
	return out
}

func (sr *seekResult) idInt(id []byte) int {
	if sr.idIdx >= len(id) {
		return -1
	}
	return int(id[sr.idIdx])
}

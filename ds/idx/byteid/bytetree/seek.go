package bytetree

import "bytes"

type seekResult struct {
	*node
	idIdx int
	found bool
}

func (bt *byteIdxByteTree) seek(id []byte) *seekResult {
	sr := &seekResult{
		node: bt.root,
	}
	for _, b := range id {
		child := sr.children[b]
		if child == nil {
			return sr
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
	n.idx = -1
	n.rest = nil
}

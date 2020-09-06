package hextree

import "bytes"

type stackFrame struct {
	b byte
	h *high
	l *low
}

type seekResult struct {
	l     *low
	idIdx int
	found bool
	stack []stackFrame
}

func (ht *hexTree) seek(id []byte, stack bool) *seekResult {
	sr := &seekResult{
		l: ht.root,
	}
	for _, b := range id {
		h := sr.l.child(b)
		l := h.child(b)
		if h == nil || l == nil {
			return sr
		}

		if stack {
			sr.stack = append(sr.stack, stackFrame{
				b: b,
				h: h,
				l: sr.l,
			})
		}
		sr.l = l
		sr.idIdx++
		sr.found = sr.match(id)
		if sr.found {
			return sr
		}
	}
	return sr
}

func (sr *seekResult) match(id []byte) bool {
	return sr.l.idx != -1 && bytes.Equal(id[sr.idIdx:], sr.l.rest)
}

func (sr *seekResult) insert(id []byte, idx int) {
	for _, b := range id[sr.idIdx:] {
		if sr.l.idx != -1 && len(sr.l.rest) > 0 {
			sr.l.bump()
		}
		h := sr.l.child(b)
		l := h.child(b)
		if h == nil || l == nil {
			sr.l.insert(b, idx, id[sr.idIdx+1:])
			return
		}
		sr.l = l
		sr.idIdx++
	}
	if sr.l.idx != -1 {
		panic("something has gone terribly wrong")
	}
	sr.l.idx = idx
}

func (l *low) bump() {
	l.insert(l.rest[0], l.idx, l.rest[1:])

	l.idx = -1
	l.rest = nil
}

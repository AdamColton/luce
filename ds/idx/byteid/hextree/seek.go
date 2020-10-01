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

func (sr *seekResult) del(id []byte) {
	sr.l.idx = -1
	sr.l.rest = nil
	// prune tree
	ln := len(sr.stack)
	for sr.l.idx == -1 && sr.l.count == 0 && ln > 0 {
		sf := sr.stack[ln-1]
		sr.idIdx--
		b := id[sr.idIdx]

		sf.h.removeChild(b)
		if sf.h.count > 0 {
			break
		}
		sf.l.removeChild(b)
		sr.l = sf.l
		ln--
	}
}

// rightThenUp searches for a non-nil child proceeding right across each
// set of children, then if none are found moving up to the parent and trying
// again.
func (sr *seekResult) rightThenUp() bool {
	l := sr.l
	for i := len(sr.stack) - 1; i >= 0; i-- {
		b := sr.stack[i].b
		for j := (b >> 4) + 1; j < 16; j++ {
			if l.children[j] != nil {
				sr.downThenLeft(j, l.children[j])
				return true
			}
		}
		h := sr.stack[i].h
		l = sr.stack[i].l
		sr.stack = sr.stack[:i]
		for j := (b & 15) + 1; j < 16; j++ {
			if h.children[j] != nil {
				sr.stack = append(sr.stack, stackFrame{
					b: (b & 240) | j,
					h: h,
					l: l,
				})
				sr.l = h.children[j]
				sr.downThenLeft(0, nil)
				return true
			}
		}
	}
	return false
}

func (sr *seekResult) downThenLeft(b byte, h *high) {
outer:
	for {
		if h != nil {
			for i := byte(0); i < 16; i++ {
				if h.children[i] != nil {
					sr.stack = append(sr.stack, stackFrame{
						b: (b << 4) | i,
						h: h,
						l: sr.l,
					})
					sr.l = h.children[i]
					h = nil
					continue outer
				}
			}
			panic("something has gone terribly wrong")
		}
		for i := byte(0); i < 16; i++ {
			if sr.l.children[i] != nil {
				h = sr.l.children[i]
				b = i
				continue outer
			}
		}
		return
	}
}

func (sr *seekResult) value() []byte {
	out := make([]byte, len(sr.stack)+len(sr.l.rest))
	for i, sf := range sr.stack {
		out[i] = sf.b
	}
	copy(out[len(sr.stack):], sr.l.rest)
	return out
}

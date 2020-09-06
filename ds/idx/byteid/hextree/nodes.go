package hextree

type high struct {
	children [16]*low
	count    byte
}

type low struct {
	children   [16]*high
	count      byte
	rest       []byte
	idx        int
	prev, next *low
}

func (l *low) child(b byte) *high {
	return l.children[b>>4]
}

func (l *low) removeChild(b byte) {
	l.children[b>>4] = nil
	l.count--
}

func (h *high) removeChild(b byte) {
	h.children[b&15] = nil
	h.count--
}

func (h *high) child(b byte) *low {
	if h == nil {
		return nil
	}
	return h.children[b&15]
}

func (l *low) insert(b byte, idx int, rest []byte) *low {
	h := l.child(b)
	if h == nil {
		h = &high{}
		l.children[b>>4] = h
	}

	l2 := &low{
		rest: make([]byte, len(rest)),
		idx:  idx,
	}
	copy(l2.rest, rest)
	l.count++
	h.count++
	h.children[b&15] = l2
	return l2
}

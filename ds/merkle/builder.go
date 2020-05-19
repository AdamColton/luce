package merkle

import (
	"hash"
)

type builder struct {
	maxSize, branch int
	h               hash.Hash
}

// NewBuilder creates a Builder whose trees will follow the limits set by
// maxSize and branch and will use the provided hash.
func NewBuilder(maxSize uint32, branch byte, h hash.Hash) Builder {
	return builder{
		maxSize: int(maxSize),
		branch:  int(branch),
		h:       h,
	}
}

func (b builder) makeNode(ln int) *branch {
	return &branch{
		children: make([]node, ln),
	}
}

// nodeLink is a link in a linked list where the payload is a node.
// start and end are the index range relative to the data being built into a
// tree.
type nodeLink struct {
	node
	next       *nodeLink
	start, end int
}

func (b builder) linker(data []byte, idx uint32, start, size, r int) *nodeLink {
	if len(data) == start {
		return nil
	}
	end := start + size
	if r > 0 {
		end++
		r--
	}
	return &nodeLink{
		node:  newDataLeaf(data[start:end], idx, b.h),
		next:  b.linker(data, idx+1, end, size, r),
		start: start,
		end:   end,
	}
}

type nodeList struct {
	ln, r     int
	root, cur *nodeLink
	data      []byte
}

func (l *nodeList) reset() {
	l.cur = l.root
	l.r = l.ln
}

func (l *nodeList) insertNodeAtCursor(b builder, ln int) *branch {
	// create node
	n := b.makeNode(ln)
	b.h.Reset()
	subcur := l.cur
	var c uint32
	d := 0
	for i := range n.children {
		b.h.Write(subcur.node.Digest())
		n.children[i] = subcur.node
		c += subcur.node.Count()
		if sd := subcur.node.Depth(); sd > d {
			d = sd
		}
		l.cur.end = subcur.end
		n.idx = subcur.maxIdx()
		subcur = subcur.next
	}
	b.h.Write(uint32ToSlice(n.idx))
	n.digest = b.h.Sum(nil)
	n.count = c
	n.depth = d + 1

	// relink
	l.cur.node = n
	l.cur.next = subcur
	n.data = l.data[l.cur.start:l.cur.end]
	l.ln -= b.branch - 1
	// move next
	l.cur = l.cur.next
	l.r -= b.branch

	return n
}

func (b builder) Build(data []byte) Tree {
	// figure out block count
	c := len(data) / b.maxSize
	if len(data)%b.maxSize != 0 {
		c++
	}

	l := &nodeList{
		ln:   c,
		root: b.linker(data, 0, 0, len(data)/c, len(data)%c),
		data: data,
	}
	l.reset()

	for l.ln > b.branch {
		if l.r < b.branch {
			l.reset()
		}
		l.insertNodeAtCursor(b, b.branch)
	}

	l.reset()
	return l.insertNodeAtCursor(b, l.ln)
}

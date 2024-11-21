package linkedlist

const empty = ^uint32(0)

type sliceRecord[T any] struct {
	prev, next uint32
	payload    T
}

type sliceNode[T any] struct {
	idx uint32
	l   *sliceList[T]
}

func (n *sliceNode[T]) Copy() Node[T] {
	cp := *n
	return &cp
}

func (n *sliceNode[T]) Prev() bool {
	n.idx = n.l.nodes[n.idx].prev
	return n.idx == empty
}

func (n *sliceNode[T]) Next() bool {
	n.idx = n.l.nodes[n.idx].next
	return n.idx == empty
}

func (n *sliceNode[T]) Get() (T, bool) {
	if n.idx == empty {
		return n.l.zero, false
	}
	return n.l.nodes[n.idx].payload, true
}

func (n *sliceNode[T]) Set(payload T) bool {
	ok := n.idx != empty
	if ok {
		n.l.nodes[n.idx].payload = payload
	}
	return ok
}

func (n *sliceNode[T]) Head() bool {
	n.idx = n.l.head
	return n.idx != empty
}

func (n *sliceNode[T]) Tail() bool {
	n.idx = n.l.tail
	return n.idx != empty
}

func (n *sliceNode[T]) Append(payload ...T) {
	last := n.l.nodes[n.idx].next
	for _, p := range payload {
		n.l.nodes[n.idx].next = n.l.getNode(p)
		n.idx = n.l.nodes[n.idx].next
	}
	n.l.nodes[n.idx].next = last
}

func (n *sliceNode[T]) Prepend(payload ...T) {
	last := n.l.nodes[n.idx].prev
	for _, p := range payload {
		n.l.nodes[n.idx].prev = n.l.getNode(p)
		n.idx = n.l.nodes[n.idx].prev
	}
	n.l.nodes[n.idx].prev = last
}

type sliceList[T any] struct {
	nodes      []sliceRecord[T]
	head, tail uint32
	freeHead   uint32
	zero       T
}

func (sl *sliceList[T]) getNode(payload T) (idx uint32) {
	if sl.freeHead != empty {
		idx = sl.freeHead
		sl.freeHead = sl.nodes[idx].next
	} else {
		idx = uint32(len(sl.nodes))
		sl.nodes = append(sl.nodes, sliceRecord[T]{
			payload: payload,
		})
	}
	return
}

package linkedlist

type poolRecord[T any] struct {
	prev, next *poolRecord[T]
	payload    T
}

type listNode[T any] struct {
	*poolRecord[T]
	*list[T]
}

type list[T any] struct {
	head, tail *poolRecord[T]
}

type unpooledlist[T any] struct {
	head, tail *poolRecord[T]
	zero       T
}

type poolList[T any] struct {
	list[T]
	p *Pool
}

type Pool[T any] struct {
	zero     T
	freeHead *poolRecord[T]
}

func NewPool[T any](zero T) *Pool[T] {
	return &Pool[T]{
		zero: zero,
	}
}

func (p *Pool[T]) NewList() List[T] {
	return &list[T]{
		p: p,
	}
}

func NewList[T any](zero T) List[T] {
	return unpooledlist[T]{
		zero: zero,
	}
}

func (n *listNode[T]) Prev() bool {
	n.poolRecord = n.prev
	return n.poolRecord != nil
}

func (n *listNode[T]) Next() bool {
	n.poolRecord = n.next
	return n.poolRecord != nil
}

func (n *listNode[T]) Get() (T, bool) {
	if n.poolRecord == nil {
		return n.list.p.zero, false
	}
	return n.payload, true
}

func (n *listNode[T]) Set(payload T) {
	n.payload = payload
}

func (l *list[T]) Append(payload T) Node[T] {
	n := l.getNode(payload)

	n.prev, n.next = l.tail, nil
	l.tail = n
	return n
}

func (l *list[T]) Prepend(payload T) Node[T] {
	n := l.getNode(payload)

	n.prev, n.next = nil, l.head
	l.head = n
	return n
}

func (l *list[T]) getNode(payload T) *poolRecord[T] {
	if l.p != nil {
		n := l.p.getNode(payload)
		return n

	}
	return &poolRecord[T]{
		payload: payload,
	}
}

func (p *Pool[T]) getNode(payload T) *poolRecord[T] {
	if p.freeHead == nil {
		return &poolRecord[T]{
			payload: payload,
		}
	}
	n := p.freeHead
	p.freeHead = n.next
	n.payload = payload
	return n
}

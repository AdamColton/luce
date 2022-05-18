package linkedlist

type record[T any] struct {
	prev, next *record[T]
	payload    T
}

type node[T any] struct {
	*record[T]
	*headTail[T]
}

type headTail[T any] struct {
	head, tail *record[T]
}

func (n *node[T]) Nil() bool {
	return n.record == nil
}

func (n *node[T]) Copy() Node[T] {
	cp := *n
	return &cp
}

func (n *node[T]) Prev() {
	if n.record != nil {
		n.record = n.prev
	}
}

func (n *node[T]) Next() {
	if n.record != nil {
		n.record = n.next
	}
}

func (n *node[T]) Get() (payload T) {
	if n.record != nil {
		payload = n.payload
	}
	return
}

func (n *node[T]) Set(payload T) {
	if n.record == nil {
		n.Append(payload)
	} else {
		n.payload = payload
	}
}

func (n *node[T]) Head() {
	n.record = n.head
}

func (n *node[T]) Tail() {
	n.record = n.tail
}

func (n *node[T]) Append(payload ...T) {
	if len(payload) == 0 {
		return
	}
	if n.record == nil {
		n.head = &record[T]{
			payload: payload[0],
		}
		n.record = n.head
		payload = payload[1:]
	}
	next := n.record.next

	for _, p := range payload {
		n.record.next = &record[T]{
			prev:    n.record,
			payload: p,
		}
		n.record = n.record.next
	}
	if next == nil {
		n.tail = n.record
	} else {
		n.record.next = next
	}
}

func (n *node[T]) Prepend(payload ...T) {
	if len(payload) == 0 {
		return
	}
	if n.record == nil {
		n.tail = &record[T]{
			payload: payload[0],
		}
		n.record = n.tail
		payload = payload[1:]
	}
	prev := n.record.next

	for _, p := range payload {
		n.record.prev = &record[T]{
			next:    n.record,
			payload: p,
		}
		n.record = n.record.prev
	}
	if prev == nil {
		n.head = n.record
	} else {
		n.record.prev = prev
	}
}

func Simple[T any]() Node[T] {
	return &node[T]{
		headTail: &headTail[T]{},
	}
}

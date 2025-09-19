package lset

// FloodProc is invoked by Flood. It is invoked once for each value of T in
// the set. Calling add checks if the value is already in the Set
type FloodProc[T comparable] func(t T, add func(T))

type stackNode[T comparable] struct {
	next *stackNode[T]
	t    T
}

func (sn *stackNode[T]) add(t T) *stackNode[T] {
	return &stackNode[T]{
		t:    t,
		next: sn,
	}
}

// Flood manages a stack of T to process. Each time FloodProc calls add, the
// value passed in will be added to the stack if it is not already in the set.
// This means that each value that FloodProc passes into add will be passed into
// FloodProc exactly once, regardless of how many times it was passed into add.
//
// Flood may not be the most descriptive name of what it is actually doing, but
// this algorithm is often used to flood a region during rasterization. Each
// cell iterates over it's neighbors calling add on each that is inside the
// region. Even if a cell is added multiple times, it will only be evaluated
// once. The resulting Set is all the cells inside the region.
func Flood[T comparable](fn FloodProc[T], init ...T) *Set[T] {
	out := New(init...)
	out.Flood(fn)
	return out
}

// Flood uses the values in the Set as the starting values for the flood
// process.
func (s *Set[T]) Flood(fn FloodProc[T]) {
	var stack *stackNode[T]
	s.All(func(t T) {
		stack = stack.add(t)
	})
	var pool *stackNode[T]
	var zero T
	add := func(t T) {
		if !s.Checksert(t) {
			if pool == nil {
				stack = stack.add(t)
			} else {
				sn := pool
				// pop a node off the pool
				pool = pool.next
				// set the value and add it to the stack
				sn.next, sn.t, stack = stack, t, sn
			}
		}
	}
	for stack != nil {
		sn, t := stack, stack.t
		// pop the stack
		stack = stack.next
		//move the unused node to the pool
		sn.t, sn.next, pool = zero, pool, sn
		fn(t, add)
	}
}

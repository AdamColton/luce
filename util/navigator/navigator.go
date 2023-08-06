package navigator

import "github.com/adamcolton/luce/ds/slice"

// Nexter represents a node structure that can get the next node based on the
// key. If "create" is true, the Next method should create the node if it does
// not exist. Context allows arbitrary additional information to be passed in
// that may be necessary to perform the Next operation.
type Nexter[Key, Node, Context any] interface {
	Next(key Key, create bool, ctx Context) (Node, bool)
}

// VoidContext is provided as a helper for instances where Context is not
// necessary.
type VoidContext struct{}

// Void is an instance of VoidContext
var Void VoidContext

// Navigator can move through graph like structures whose nodes fulfill Nexter.
type Navigator[Key any, Node Nexter[Key, Node, Context], Context any] struct {
	// Cur is current or cursor node.
	Cur Node

	// Idx of the next Key to navigate
	Idx int
	// Keys to Navigate
	Keys slice.Slice[Key]

	// TraceNodes causes Nodes to be populated during a Seek operation if true.
	TraceNodes bool
	// Nodes iterated over during a seek operation.
	Nodes slice.Slice[Node]
}

// Trace is a chainable helper. It sets TraceNodes to trace and returns the
// Navigator.
func (n *Navigator[Key, Node, Context]) Trace(trace bool) *Navigator[Key, Node, Context] {
	n.TraceNodes = trace
	return n
}

// Seek starts with the initial value of Cur and iterates through Keys calling
// Next for each key so long as the returned OK is true. If TraceNodes is
// true, the navigated nodes will be added to Nodes.
func (n *Navigator[Key, Node, Context]) Seek(create bool, ctx Context) (next Node, ok bool) {
	next, ok = n.Cur, true
	for ; n.Idx < len(n.Keys); n.Idx++ {
		k := n.IdxKey()
		next, ok = n.Cur.Next(k, create, ctx)
		if n.TraceNodes {
			n.Nodes = append(n.Nodes, n.Cur)
		}
		n.Cur = next
		if !ok {
			break
		}
	}
	return
}

// IdxKey is a helper to get the key at the current index
func (n *Navigator[Key, Node, Context]) IdxKey() Key {
	return n.Keys[n.Idx]
}

// Pop the last value in Nodes off and assign it to Cur. This requires that
// Nodes are populated. So in order to use this after a Seek operation,
// TraceNodes needs to be true prior to the seek operation.
func (n *Navigator[Key, Node, Context]) Pop() (node Node, ok bool) {
	ok = len(n.Nodes) > 0
	if ok {
		node = n.Cur
		n.Cur, n.Nodes = n.Nodes.Pop()
		n.Idx--
	}
	return
}

// Package navigator moves through graph-like structures by following a sequence
// of keys.
package navigator

import "github.com/adamcolton/luce/ds/slice"

// Nexter represents a node structure that can get the next node based on the
// key. If "create" is true, the Next method should create the node if it does
// not exist. Context allows arbitrary additional information to be passed in
// that may be necessary to perform the Next operation. The 'ok' bool indicates
// if operation was successful and the navigator should continue.
type Nexter[Key, Node, Context any] interface {
	Next(key Key, create bool, ctx Context) (node Node, ok bool)
}

// VoidContext is provided as a helper for instances where Context is not
// necessary.
type VoidContext struct{}

// Void is an instance of VoidContext
var Void VoidContext

// Navigator can move through graph like structures whose nodes fulfill Nexter.
// The Node type is its own Nexter, so following a key leads to another Node.
type Navigator[Key any, Node Nexter[Key, Node, Context], Context any] struct {
	// Cur is current or cursor node. Seek starts here.
	Cur Node

	// Idx of the next Key to navigate. Seek starts here.
	Idx int
	// Keys to Navigate
	Keys slice.Slice[Key]

	// TraceNodes causes Nodes to be populated during a Seek operation if true.
	TraceNodes bool
	// Nodes iterated over during a seek operation.
	Nodes slice.Slice[Node]
}

// New creates a Navigator that starts at start and will follow keys. TraceNodes
// is false.
func New[Key any, Node Nexter[Key, Node, Context], Context any](start Node, keys []Key) *Navigator[Key, Node, Context] {
	return &Navigator[Key, Node, Context]{
		Cur:  start,
		Keys: keys,
	}
}

// Trace is a chainable helper. It sets TraceNodes to trace and returns the
// Navigator.
func (n *Navigator[Key, Node, Context]) Trace(trace bool) *Navigator[Key, Node, Context] {
	n.TraceNodes = trace
	return n
}

// Seek starts with the current value of Cur and iterates through Keys, from
// Idx, calling Next for each key so long as the returned OK is true. Cur and Idx
// are updated as it goes, so a Seek can be continued.
//
// If Next returns false, Seek stops and returns that node and false. Cur is
// then that node (generally the zero value) and Idx is the index of the key that
// failed. If there are no keys left it returns Cur and true.
//
// If TraceNodes is true, the node each key was followed from is added to Nodes.
// So after a Seek that follows every key, Nodes holds every node visited except
// the last, which is Cur.
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

// IdxKey is a helper to get the key at the current index. It panics if Idx is
// not in the range of Keys.
func (n *Navigator[Key, Node, Context]) IdxKey() Key {
	return n.Keys[n.Idx]
}

// Pop the last value in Nodes off and assign it to Cur, and step Idx back. This
// requires that Nodes are populated. So in order to use this after a Seek
// operation, TraceNodes needs to be true prior to the seek operation. The node
// that was Cur is returned. If Nodes is empty, it returns the zero value and
// false and changes nothing.
//
// It undoes the last step of a Seek that followed every key. After a Seek that
// failed, Idx is one key too far back, because it was not advanced for the key
// that failed.
func (n *Navigator[Key, Node, Context]) Pop() (node Node, ok bool) {
	ok = len(n.Nodes) > 0
	if ok {
		node = n.Cur
		n.Cur, n.Nodes = n.Nodes.Pop()
		n.Idx--
	}
	return
}

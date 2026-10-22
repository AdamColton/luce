package depgraph

import (
	"reflect"
	"slices"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
)

// Node is any node in a graph, either a Definition or a Value, without its
// type parameter.
type Node interface {
	Name() string
	Type() reflect.Type
	Uses(in *Context) list.Wrapper[Node]
	UsedBy(in *Context) list.Wrapper[Node]
	Value(in *Context) any
}

// getter is a Node that can produce its value as T, used to resolve a node
// found by type.
type getter[T any] interface {
	get(in *Context) T
}

// edges is the part of a node's state in a Context that doesn't depend on its
// type.
type edges struct {
	// run is 0 when the value needs computing, 1 while it is being computed
	// (so re-entry is a cycle) and 2 once it is done. Values don't use it.
	run    byte
	uses   slice.Slice[Node]
	usedBy slice.Slice[Node]
}

func (e *edges) base() *edges {
	return e
}

// state is a node's state in a Context.
type state[T any] struct {
	edges
	val T
}

type stater interface {
	base() *edges
}

// stateOf returns n's state in the Context, creating it if create is set. It
// returns nil when n has no state there and create is not set.
func stateOf[T any](n Node, in *Context, create bool) *state[T] {
	s, found := in.nodes[n]
	if !found {
		if !create {
			return nil
		}
		s = &state[T]{}
		in.nodes[n] = s
	}
	return s.(*state[T])
}

// edgesOf returns n's edges in the Context, or nil if n has no state there.
func (c *Context) edgesOf(n Node) *edges {
	if s, found := c.nodes[n]; found {
		return s.base()
	}
	return nil
}

// used records that the node being computed, if any, used n.
func (c *Context) used(n Node, e *edges) {
	if c.cur == nil || slices.Contains(e.usedBy, c.cur) {
		return
	}
	e.usedBy = append(e.usedBy, c.cur)
	cur := c.edgesOf(c.cur)
	cur.uses = append(cur.uses, n)
}

// clearUses drops the edges from n's last computation, so that after it
// recomputes, its edges are only the dependencies that run actually used.
func (c *Context) clearUses(n Node, e *edges) {
	for _, use := range e.uses {
		ue := c.edgesOf(use)
		if idx := slices.Index(ue.usedBy, n); idx >= 0 {
			ue.usedBy = ue.usedBy.RemoveOrdered(idx)
		}
	}
	e.uses = e.uses[:0]
}

func (c *Context) uses(n Node) list.Wrapper[Node] {
	var uses slice.Slice[Node]
	if e := c.edgesOf(n); e != nil {
		uses = e.uses
	}
	return list.Wrap(uses)
}

func (c *Context) usedBy(n Node) list.Wrapper[Node] {
	var usedBy slice.Slice[Node]
	if e := c.edgesOf(n); e != nil {
		usedBy = e.usedBy
	}
	return list.Wrap(usedBy)
}

func (c *Context) allUsedBy(n Node) *lset.Set[Node] {
	fn := func(n Node, add func(Node)) {
		c.usedBy(n).For(add)
	}
	return lset.Flood(fn, n)
}

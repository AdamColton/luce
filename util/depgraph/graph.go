package depgraph

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
)

// Context is a concrete instance of a graph. It holds every node's value and
// edges for that instance.
type Context struct {
	cur    Node
	nodes  map[Node]stater
	byName lmap.Wrapper[string, Node]
	byType lmap.Wrapper[reflect.Type, Node]
}

func (c *Context) Each(fn func(name string, n Node, done *bool)) {
	c.byName.Each(func(name string, n Node, done *bool) {
		fn(name, n, done)
	})
}

func (c *Context) Resolve(name string) Node {
	return c.byName.GetVal(name)
}

const ErrNotFound = lerr.Str("not found")

// Resolve returns the value of the node added to the Context by name, either a
// Value or a Definition. It panics with ErrNotFound if there is none, and with
// lerr.ErrTypeMismatch if the node's type is not T.
func Resolve[T any](name string, in *Context) T {
	n, found := in.byName.Get(name)
	if !found {
		panic(ErrNotFound)
	}
	g, ok := n.(getter[T])
	if !ok {
		panic(lerr.ErrTypeMismatch{
			Expected: reflect.TypeFor[T](),
			Actual:   n.Type(),
		})
	}
	return g.get(in)
}

// ResolveType returns the value of the node added to the Context for type T,
// either a Value or a Definition. It panics with ErrNotFound if there is none.
func ResolveType[T any](in *Context) T {
	n, found := in.byType.Get(reflect.TypeFor[T]())
	if !found {
		panic(ErrNotFound)
	}
	return n.(getter[T]).get(in)
}

// Set sets value to a copy of what depgraph.ResolveType[*V] returns.
func Set[V any](value *V, in *Context) {
	*value = *ResolveType[*V](in)
}

func NewContext() *Context {
	return &Context{
		nodes:  make(map[Node]stater),
		byType: lmap.New[reflect.Type, Node](nil),
		byName: lmap.New[string, Node](nil),
	}
}

const (
	ErrCycle          = lerr.Str("graph cycles")
	ErrDoubleRegister = lerr.Str("cannot register a type more than once")
	ErrNotRegistered  = lerr.Str("type not registered")
)

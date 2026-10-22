package depgraph

import (
	"reflect"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/lerr"
)

// Definition is an abstract node in a graph: a name and a generator that
// computes the node's value from the Context. It holds no values itself; each
// Context it is resolved in keeps its own.
type Definition[T any] struct {
	name string
	fn   func(*Context) T
}

const ErrNameCollision = lerr.Str("name collision")

// Define creates a Definition. It is added to a Context's names the first time
// it is resolved there.
func Define[T any](name string, fn func(in *Context) T) *Definition[T] {
	return &Definition[T]{
		fn:   fn,
		name: name,
	}
}

// DefineType creates a Definition and adds it to the Context by name and by
// type, so depgraph.ResolveType[T] finds it. It replaces anything already there
// for T.
func DefineType[T any](name string, fn func(in *Context) T, in *Context) *Definition[T] {
	d := Define(name, fn)
	in.byName.CheckSet(name, d)
	in.byType.Set(reflect.TypeFor[T](), d)
	return d
}

// Resolve returns the node's value in the Context, computing it if it hasn't
// been computed or something it depends on has changed. When called from a
// generator, it records that the generator's node uses this one. It panics with
// ErrCycle if the node depends on itself.
func (d *Definition[T]) Resolve(in *Context) T {
	s := stateOf[T](d, in, true)
	if s.run == 1 {
		panic(ErrCycle)
	}
	if s.run == 0 {
		prev := in.cur
		in.cur = d
		s.run = 1

		lmap.CheckSet(in.byName, d.name, Node(d))
		in.clearUses(d, &s.edges)

		s.val = d.fn(in)
		s.run = 2
		in.cur = prev
	}
	in.used(d, &s.edges)
	return s.val
}

func (d *Definition[T]) get(in *Context) T {
	return d.Resolve(in)
}

func (d *Definition[T]) Name() string {
	return d.name
}

func (d *Definition[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

// Uses returns the nodes this node used when it was resolved in the Context.
// It is empty if the node has not been resolved there.
func (d *Definition[T]) Uses(in *Context) list.Wrapper[Node] {
	return in.uses(d)
}

// UsedBy returns the nodes that used this node when they were resolved in the
// Context. It is empty if the node has not been resolved there.
func (d *Definition[T]) UsedBy(in *Context) list.Wrapper[Node] {
	return in.usedBy(d)
}

// Value returns the node's value in the Context, resolving it first if it
// hasn't been.
func (d *Definition[T]) Value(in *Context) any {
	return d.Resolve(in)
}

// AllUsedBy returns the node and every node that depends on it, directly or
// indirectly, in the Context.
func (d *Definition[T]) AllUsedBy(in *Context) *lset.Set[Node] {
	return in.allUsedBy(d)
}

package depgraph

import (
	"reflect"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/lerr"
)

// Value is a node whose value is set rather than computed. It belongs to the
// Context it was designated in.
type Value[T any] struct {
	name string
	ctx  *Context
}

// Designate adds a Value to the Context by name.
func Designate[T any](name string, value T, in *Context) *Value[T] {
	v := &Value[T]{
		name: name,
		ctx:  in,
	}
	stateOf[T](v, in, true).val = value
	in.byName.CheckSet(name, v)
	return v
}

// DesignateType adds a Value to the Context by name and by type, so
// depgraph.ResolveType[T] finds it. It replaces anything already there for T.
func DesignateType[T any](name string, value T, in *Context) *Value[T] {
	v := Designate(name, value, in)
	in.byType.Set(reflect.TypeFor[T](), v)
	return v
}

type Named interface {
	Name() string
}

// DesignateNamedType calls DesignateType with the value's own name.
func DesignateNamedType[V Named](value V, in *Context) *Value[V] {
	return DesignateType(value.Name(), value, in)
}

// Ctx returns the Context the Value was designated in.
func (v *Value[T]) Ctx() *Context {
	return v.ctx
}

// Get returns the Value. When called from a generator, it records that the
// generator's node uses this one.
func (v *Value[T]) Get() T {
	s := stateOf[T](v, v.ctx, false)
	v.ctx.used(v, &s.edges)
	return s.val
}

func (v *Value[T]) get(*Context) T {
	return v.Get()
}

// Update sets the Value. Every node that depends on it is recomputed the next
// time it is resolved.
func (v *Value[T]) Update(value T) {
	self := Node(v)
	v.AllUsedBy().Each(func(n Node, done *bool) {
		if n != self {
			v.ctx.edgesOf(n).run = 0
		}
	})
	stateOf[T](v, v.ctx, false).val = value
}

// UpdateAny calls Update with value, for a caller that only knows the Value's
// type at runtime. It returns lerr.ErrTypeMismatch if value is not a T.
func (v *Value[T]) UpdateAny(value any) error {
	t, ok := value.(T)
	if !ok {
		return lerr.ErrTypeMismatch{
			Expected: reflect.TypeFor[T](),
			Actual:   reflect.TypeOf(value),
		}
	}
	v.Update(t)
	return nil
}

func (v *Value[T]) Name() string {
	return v.name
}

func (v *Value[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

// Uses is always empty; a Value doesn't depend on other nodes.
func (v *Value[T]) Uses(in *Context) list.Wrapper[Node] {
	return in.uses(v)
}

// UsedBy returns the nodes that used this Value when they were resolved in the
// Context. It is empty in any Context but the Value's own.
func (v *Value[T]) UsedBy(in *Context) list.Wrapper[Node] {
	return in.usedBy(v)
}

// Value returns the Value as an any if in is the Value's Context, and nil
// otherwise.
func (v *Value[T]) Value(in *Context) any {
	if in != v.ctx {
		return nil
	}
	return v.Get()
}

// AllUsedBy returns the Value and every node that depends on it, directly or
// indirectly.
func (v *Value[T]) AllUsedBy() *lset.Set[Node] {
	return v.ctx.allUsedBy(v)
}

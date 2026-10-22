package depgraph

import "reflect"

// DefaultValue is a Value's name and default. It isn't bound to a Context; a
// Context that hasn't designated the name gets the default when it is
// resolved there.
type DefaultValue[T any] struct {
	name  string
	value T
}

// TypeRegistrar reads Values from a config, and can designate a new Value of
// any type registered with it. *depcfg.Parser is one.
type TypeRegistrar interface {
	// RegisterType lets the TypeRegistrar designate a new Value of type t, by
	// calling designate with a value of that type.
	RegisterType(t reflect.Type, designate func(name string, v any, in *Context))
}

// DesignateDefault creates a DefaultValue, and registers T with each TypeRegistrar so a
// config can designate a new T. If T is a pointer, every Context that falls
// back to the default shares the value it points to.
func DesignateDefault[T any](name string, v T, rs ...TypeRegistrar) *DefaultValue[T] {
	for _, r := range rs {
		r.RegisterType(reflect.TypeFor[T](), func(name string, v any, in *Context) {
			Designate(name, v.(T), in)
		})
	}
	return &DefaultValue[T]{
		name:  name,
		value: v,
	}
}

// Resolve returns the value of the node added to the Context by the
// DefaultValue's name. If there is none, it first designates the default
// there. It panics with lerr.ErrTypeMismatch if the node's type is not T.
func (d *DefaultValue[T]) Resolve(in *Context) T {
	if _, found := in.byName.Get(d.name); !found {
		Designate(d.name, d.value, in)
	}
	return Resolve[T](d.name, in)
}

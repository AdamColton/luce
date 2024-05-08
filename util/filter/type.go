package filter

import (
	"reflect"

	"github.com/adamcolton/luce/math/ints"
	"github.com/adamcolton/luce/util/reflector"
)

// AnyType is a type filter that returns True for any type. This can be useful
// when creating checks for things like functions.
func AnyType() Type {
	return Type{func(t reflect.Type) bool { return true }}
}

// Type is a wrapper around Filter[reflect.Type] to provide helper logic
// for type filtering.
type Type struct {
	Filter[reflect.Type]
}

// Func builds a filter for a function. Both args and rets can be either
// a filter.Type, filter.Filter[reflect.Type] or reflect.Type.
func Func(args, rets []any) Type {
	// TODO: panic on bad type.
	af := make([]Filter[reflect.Type], len(args))
	rf := make([]Filter[reflect.Type], len(rets))
	for i, a := range args {
		switch t := a.(type) {
		case reflect.Type:
			af[i] = IsType(t).Filter
		case Type:
			af[i] = t.Filter
		case Filter[reflect.Type]:
			af[i] = t
		}
	}
	for i, r := range rets {
		switch t := r.(type) {
		case reflect.Type:
			rf[i] = IsType(t).Filter
		case Type:
			rf[i] = t.Filter
		case Filter[reflect.Type]:
			rf[i] = t
		}
	}
	return Type{funcFilter(af, rf)}
}

func funcFilter(args, rets []Filter[reflect.Type]) Filter[reflect.Type] {
	return func(t reflect.Type) bool {
		if t.Kind() != reflect.Func || t.NumIn() != len(args) || t.NumOut() != len(rets) {
			return false
		}
		for i, arg := range args {
			if !arg(t.In(i)) {
				return false
			}
		}
		for i, ret := range rets {
			if !ret(t.Out(i)) {
				return false
			}
		}
		return true
	}
}

// NumInEq is a helper equivalent to filter.NumIn(filter.EQ(n)).
func NumInEq(n int) Type {
	return NumIn(EQ(n))
}

// NumOutEq is a helper equivalent to filter.NumOut(filter.EQ(n)).
func NumOutEq(n int) Type {
	return NumOut(EQ(n))
}

// InType is a helper equivalent to filter.IsType(t).In(n).
func InType(n int, t reflect.Type) Type {
	return IsType(t).In(n)
}

// OutType is a helper equivalent to filter.IsType(t).Out(n).
func OutType(n int, t reflect.Type) Type {
	return IsType(t).Out(n)
}

// Implements returns true if the provided type implements interface I.
func Implements[I any]() Type {
	i := reflector.Type[I]()
	return Type{func(t reflect.Type) bool {
		return t.Implements(i)
	}}
}

// OnInterface applies the filter to the TypeOf i.
func (t Type) OnInterface(i any) bool {
	return t.Filter(reflect.TypeOf(i))
}

// Elem checks the filter type's Elem against the underlying filter.
func (t Type) Elem() Type {
	return Type{func(t2 reflect.Type) (out bool) {
		e, ok := reflector.Elem(t2)
		return ok && t.Filter(e)
	}}
}

// In checks the filter type's agument number i against the given filter. If i
// is less than 0, it will be indexed relative to the number of arguments, so -1
// will return the last argument.
func (t Type) In(i int) Type {
	return Type{func(t2 reflect.Type) bool {
		idx, inRange := ints.Idx(i, t2.NumIn())
		return inRange && t2 != nil && t2.Kind() == reflect.Func && t.Filter(t2.In(idx))

	}}
}

// Out checks the filter type's agument number i against the given filter. If i
// is less than 0, it will be indexed relative to the number of returns, so -1
// will return the last return.
func (t Type) Out(i int) Type {
	return Type{func(t2 reflect.Type) bool {
		idx, inRange := ints.Idx(i, t2.NumOut())
		return inRange && t2 != nil && t2.Kind() == reflect.Func && t.Filter(t2.Out(idx))
	}}
}

// Method applies the Type filter to the method function.
func (t Type) Method() Filter[*reflector.Method] {
	return func(m *reflector.Method) bool {
		return t.Filter(m.Func.Type())
	}
}

// And builds a new Type filter that will return true if both underlying
// Type filters are true.
func (t Type) And(t2 Type) Type {
	return Type{t.Filter.And(t2.Filter)}
}

// Or builds a new Type filter that will return true if either underlying
// Type filters is true.
func (t Type) Or(t2 Type) Type {
	return Type{t.Filter.Or(t2.Filter)}
}

// Or builds a new Type filter that will return true if the underlying Type
// filter is false.
func (t Type) Not() Type {
	return Type{t.Filter.Not()}
}

// TypeChecker checks a value's type against a filter. It returns the underlying
// type. It returns an error if the type fails the underlying filter.
type TypeChecker func(i any) (reflect.Type, error)

// Check creates a TypeChecker from a Type filter. It uses reflector.ToType,
// so that it can accept either a reflect.Type and use it directly or an
// interface which it will call reflect.ToType on.
func (t Type) Check(errFn func(reflect.Type) error) TypeChecker {
	return func(i any) (reflect.Type, error) {
		it := reflector.ToType(i)
		if t.Filter(it) {
			return it, nil
		}
		return it, errFn(it)
	}
}

// IsKind creates a Type filter that returns true when given a type that
// matches the specified kind.
func IsKind(kind reflect.Kind) Type {
	return Type{func(t reflect.Type) bool {
		return t != nil && t.Kind() == kind
	}}
}

// IsType creates a filter using referenceType. Returns true if the filterType
// is the same as referenceType.
func IsType(referenceType reflect.Type) Type {
	return Type{func(filterType reflect.Type) bool {
		return referenceType == filterType
	}}
}

// NumIn checks the filter type's NumIn value against the given filter.
func NumIn(f func(int) bool) Type {
	return Type{func(t reflect.Type) (out bool) {
		return t != nil && t.Kind() == reflect.Func && f(t.NumIn())
	}}
}

// NumOut checks the filter type's NumOut value against the given filter.
func NumOut(f func(int) bool) Type {
	return Type{func(t reflect.Type) (out bool) {
		return t != nil && t.Kind() == reflect.Func && f(t.NumOut())
	}}
}

// MethodName takes a string filter and applies it to a Method name.
func MethodName(f func(string) bool) Filter[*reflector.Method] {
	return func(m *reflector.Method) bool {
		return f(m.Name)
	}
}

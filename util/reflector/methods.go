package reflector

import (
	"reflect"
)

// Method captures the reflection values that are useful when doing reflection
// on a method.
type Method struct {
	// On is the type the Method is define on
	On reflect.Value
	// Method is the reflect.Method value for the method
	reflect.Method
	// Func value of the method.
	Func reflect.Value
}

// MethodOn get a Method by name. The 'on' argument can be either an interface
// or a reflect.Value.
func MethodOn(on any, name string) *Method {
	v := ToValue(on)
	m, ok := v.Type().MethodByName(name)
	if !ok {
		return nil
	}
	return &Method{
		On:     v,
		Method: m,
		Func:   v.MethodByName(name),
	}
}

// AssignTo attempts to assign this method to fnPtr. The value of success
// indicates if it worked.
func (m *Method) AssignTo(fnPtr any) (success bool) {
	defer func() {
		success = recover() == nil
	}()
	reflect.ValueOf(fnPtr).Elem().Set(m.Func)
	return
}

// Methods on a single value - at least that's the intention.
type Methods []*Method

// MethodsOn returns all the methods on the interface provided.
func MethodsOn(i any) Methods {
	on := reflect.ValueOf(i)
	t := on.Type()
	ln := on.NumMethod()
	out := make(Methods, ln)
	for i := range out {
		out[i] = &Method{
			On:     on,
			Method: t.Method(i),
			Func:   on.Method(i),
		}
	}
	return out
}

// Funcs returns just the Func values of the Methods. This can be useful because
// the On and Method fields are generally only useful for filtering and the Func
// value can be invoked with .Call - so after filtering, this allows the Func
// values to be used.
func (ms Methods) Funcs() []any {
	out := make([]any, len(ms))
	for i, m := range ms {
		out[i] = m.Func.Interface()
	}
	return out
}

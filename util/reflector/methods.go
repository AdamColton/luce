package reflector

import (
	"reflect"
)

// Method captures the reflection values that are useful when doing reflection
// on a method.
type Method struct {
	On reflect.Value
	reflect.Method
	Func reflect.Value
}

func NewMethod(on reflect.Value, idx int) *Method {
	if idx < 0 || idx >= on.NumMethod() {
		return nil
	}
	return &Method{
		On:     on,
		Method: on.Type().Method(idx),
		Func:   on.Method(idx),
	}
}

func (m *Method) SetTo(i interface{}) (success bool) {
	defer func() {
		success = recover() == nil
	}()
	reflect.ValueOf(i).Elem().Set(m.Func)
	return
}

type Methods []*Method

func MethodsOn(i interface{}) Methods {
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

func (ms Methods) Funcs() []interface{} {
	out := make([]interface{}, len(ms))
	for i, m := range ms {
		out[i] = m.Func.Interface()
	}
	return out
}

type MethodFilter func(m *Method) bool

func (mf MethodFilter) On(i interface{}) Methods {
	ms := MethodsOn(i)
	rm := 0
	for i, m := range ms {
		if mf(m) {
			if rm > 0 {
				ms[i-rm] = m
			}
		} else {
			rm++
		}
	}
	return ms[:len(ms)-rm]
}

func (mf MethodFilter) One(i interface{}) *Method {
	on := reflect.ValueOf(i)
	ln := on.NumMethod()
	for i := 0; i < ln; i++ {
		if m := NewMethod(on, i); mf(m) {
			return m
		}
	}
	return nil
}

// Or builds a new MethodFilter that will return true if either underlying
// MethodFilter is true.
func (mf MethodFilter) Or(mf2 MethodFilter) MethodFilter {
	return func(val *Method) bool {
		return mf(val) || mf2(val)
	}
}

// And builds a new MethodFilter that will return true if both underlying
// MethodFilters are true.
func (mf MethodFilter) And(mf2 MethodFilter) MethodFilter {
	return func(val *Method) bool {
		return mf(val) && mf2(val)
	}
}

// Not builds a new MethodFilter that will return true if the underlying
// MethodFilter is false.
func (mf MethodFilter) Not() MethodFilter {
	return func(val *Method) bool {
		return !mf(val)
	}
}

func MethodName(f func(string) bool) MethodFilter {
	return func(m *Method) bool {
		return f(m.Name)
	}
}

func (ff FuncFilter) MethodFilter() MethodFilter {
	return func(m *Method) bool {
		return ff(m.Func.Type())
	}
}

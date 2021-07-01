package reflector

import "reflect"

type FuncFilter func(fn reflect.Type) bool

func FuncOf(i interface{}) FuncFilter {
	t := reflect.TypeOf(i)
	k := t.Kind()
	if k != reflect.Func {
		return nil
	}
	return func(fn reflect.Type) bool {
		return fn == t
	}
}

// ArgType checks that argument n is of type t. If n <0, it will match from the
// end of the args (i.e. -1 is the last argument).
func ArgType(t reflect.Type, n int) FuncFilter {
	return func(t reflect.Type) bool {
		idx := wrap(n, t.NumIn())
		if idx == -1 {
			return false
		}
		return t.In(idx) == t
	}
}

func ArgCount(f func(int) bool) FuncFilter {
	return func(t reflect.Type) bool {
		return f(t.NumIn())
	}
}

// ReturnType checks that return value n is of type t. If n <0, it will match
// from the end of the args (i.e. -1 is the last argument).
func ReturnType(t reflect.Type, n int) FuncFilter {
	return func(t reflect.Type) bool {
		idx := wrap(n, t.NumOut())
		if idx == -1 {
			return false
		}
		return t.Out(n) == t
	}
}

func ReturnCount(f func(int) bool) FuncFilter {
	return func(t reflect.Type) bool {
		return f(t.NumIn())
	}
}

var ErrFunc = ReturnType(reflect.TypeOf((error)(nil)), -1)

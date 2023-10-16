package handler

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

type retFlag int8

const (
	hasRet retFlag = 1
	hasErr retFlag = 2
)

func (rf retFlag) ErrIdx() int {
	if rf&hasErr == 0 {
		return -1
	}
	return int(rf & hasRet)
}

func (rf retFlag) HasRet() bool {
	return rf&hasRet != 0
}

// Handler is a func that takes 0 or 1 arguments, may have 1 return of any
// type and may have and error return argument. For example, if the return
// type is a string, the valid returns are (),(string),(error), (string, error).
// A Handler can also be identified by name.
type Handler struct {
	retFlag
	fn reflect.Value
}

// Type returns the type of the function argument. If the function takes no
// arguments, nil is returned.
func (h *Handler) Type() reflect.Type {
	t := h.fn.Type()
	if t.NumIn() == 0 {
		return nil
	}
	return t.In(0)
}

var nilValue reflect.Value

// Handle calls the underlying function with 'i' as the argument. If 'i' is not
// a valid argument Handle will panic.
func (h *Handler) Handle(i any) (any, error) {
	return h.HandleValue(reflect.ValueOf(i))
}

// HandleValue calls the underlying function with 'v' as the argument. If 'v' is
// not a valid argument HandleValue will panic.
func (h *Handler) HandleValue(v reflect.Value) (any, error) {
	var err error
	var args []reflect.Value
	if v != nilValue {
		args = []reflect.Value{v}
	}
	rets := h.fn.Call(args)
	if eIdx := h.ErrIdx(); eIdx >= 0 {
		if iErr := rets[eIdx].Interface(); iErr != nil {
			err = iErr.(error)
		}
	}
	var out any
	if h.HasRet() {
		out = rets[0].Interface()
	}
	return out, err
}

// New creates a Handler from a function. The 'name' argument can be blank.
func New(i any) (h *Handler, err error) {
	return ByValue(reflector.ToValue(i))
}

// ByValue creates a handler from a reflect.Value. This can be useful is the
// reflect.Value is already available.

func ByValue(v reflect.Value) (h *Handler, err error) {
	// TODO: only accept zero args if name != ""
	t := v.Type()
	if t.Kind() != reflect.Func {
		return nil, lerr.Str("expected func, got: " + t.String())
	}
	nOut := t.NumOut()
	if t.NumIn() > 1 || nOut > 2 || (nOut == 2 && t.Out(1) != ltype.Err) {
		return nil, lerr.Str("expected func(T?) (U?, error?) where ? indicates optional, got: " + t.String())
	}
	h = &Handler{
		fn: v,
	}
	if nOut > 0 && t.Out(nOut-1) == ltype.Err {
		h.retFlag = hasErr
	}
	if nOut == 2 || (nOut == 1 && h.retFlag != hasErr) {
		h.retFlag |= hasRet
	}
	return
}

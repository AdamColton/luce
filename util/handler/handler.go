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

type Handler struct {
	retFlag
	fn   reflect.Value
	name string
}

func (h *Handler) Name() string {
	return h.name
}

func (h *Handler) key() (k key) {
	t := h.fn.Type()
	if t.NumIn() == 0 {
		k.name = h.name
	} else {
		k.Type = t.In(0)
	}
	return
}

func (h *Handler) Type() reflect.Type {
	t := h.fn.Type()
	if t.NumIn() == 0 {
		return nil
	}
	return t.In(0)
}

var nilValue reflect.Value

func (h *Handler) Handle(i any) (any, error) {
	return h.HandleValue(reflect.ValueOf(i))
}

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

func New(i any, name string) (h *Handler, err error) {
	return ByValue(reflector.ToValue(i), name)
}

// TODO: only accept zero args if name != ""
func ByValue(v reflect.Value, name string) (h *Handler, err error) {
	t := v.Type()
	if t.Kind() != reflect.Func {
		return nil, lerr.Str("expected func, got: " + t.String())
	}
	nOut := t.NumOut()
	if t.NumIn() > 1 || nOut > 2 || (nOut == 2 && t.Out(1) != ltype.Err) {
		return nil, lerr.Str("expected func(T?) (U?, error?) where ? indicates optional, got: " + t.String())
	}
	h = &Handler{
		fn:   v,
		name: name,
	}
	if nOut > 0 && t.Out(nOut-1) == ltype.Err {
		h.retFlag = hasErr
	}
	if nOut == 2 || (nOut == 1 && h.retFlag != hasErr) {
		h.retFlag |= hasRet
	}
	return
}

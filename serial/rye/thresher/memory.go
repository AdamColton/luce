package thresher

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"runtime"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

var byPtr = lmap.Map[uintptr, *rootObj]{}
var byID = lmap.Map[string, *rootObj]{}

type rootObj struct {
	addr     uintptr
	v        reflect.Value
	id       []byte
	callback func(any)
}

func (ro *rootObj) init() {
	ptr := ro.addr
	byPtr[ptr] = ro
	idStr := string(ro.id)
	byID[idStr] = ro
	runtime.SetFinalizer(ro.v.Interface(), func(any) {
		delete(byPtr, ptr)
		delete(byID, idStr)
	})
}

var zeroByte = []byte{0}

func (ro *rootObj) getID() []byte {
	if ro == nil {
		return zeroByte
	}
	return ro.id
}

func newRootObj(ptr uintptr, v reflect.Value) *rootObj {
	if !reflector.CanElem(v.Kind()) {
		panic("root obj must have elem")
	}
	id := make([]byte, 32)
	lerr.Must(rand.Read(id))

	ro := &rootObj{
		addr: ptr,
		v:    v,
		id:   id,
	}

	ro.init()
	return ro
}

func rootObjByV(v reflect.Value) *rootObj {
	if k := v.Kind(); k != reflect.Pointer && k != reflect.Slice && k != reflect.Map {
		return nil
	}
	ptr := uintptr(v.UnsafePointer())
	if ptr == 0 {
		return nil
	}
	ro, found := byPtr[ptr]
	if found {
		return ro
	}
	return newRootObj(ptr, v)
}

func rootObjByID(id []byte) *rootObj {
	return byID[string(id)]
}

func awaitRootObjByID(id []byte, callback func(any)) {
	if bytes.Equal(id, zeroByte) {
		callback(nil)
		return
	}
	ro, found := byID[string(id)]
	if !found {
		ro = &rootObj{
			id: id,
		}
		byID[string(id)] = ro
	}
	if ro.v.Kind() != reflect.Invalid {
		callback(ro.v.Interface())
	} else if ro.callback == nil {
		ro.callback = callback
	} else {
		prev := ro.callback
		ro.callback = func(a any) {
			prev(a)
			callback(a)
		}
	}
}

func makeRootObj(t reflect.Type, id []byte) *rootObj {
	ro, found := byID[string(id)]
	if !found {
		ro = &rootObj{
			id: id,
		}
	}
	if ro.v.Kind() == reflect.Invalid {
		ro.v = reflect.New(t)
		ro.addr = uintptr(ro.v.UnsafePointer())
		ro.init()
	}
	if ro.callback != nil {
		ro.callback(ro.v.Interface())
		ro.callback = nil
	}

	return ro
}

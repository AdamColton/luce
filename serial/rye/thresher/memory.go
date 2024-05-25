package thresher

import (
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
	addr uintptr
	v    reflect.Value
	id   []byte
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

	byPtr[ptr] = ro
	idStr := string(id)
	byID[idStr] = ro
	runtime.SetFinalizer(v.Interface(), func(any) {
		delete(byPtr, ptr)
		delete(byID, idStr)
	})
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

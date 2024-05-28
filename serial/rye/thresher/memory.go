package thresher

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"runtime"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
)

var byPtr = lmap.Map[uintptr, *rootObj]{}
var byID = lmap.Map[string, *rootObj]{}

type storeRecord struct {
	data []byte
	t    reflect.Type
}

var store = lmap.Map[string, storeRecord]{}

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

func getStoreByID(id []byte) *rootObj {
	if bytes.Equal(id, zeroByte) {
		return nil
	}

	idStr := string(id)
	ro, found := byID[idStr]
	if !found {
		ro = &rootObj{
			id: id,
		}
	} else if ro.v.Kind() != reflect.Invalid {
		return ro
	}

	rec, found := store[idStr]
	if !found {
		return ro
	}

	ro.v = reflect.New(rec.t)
	ro.addr = uintptr(ro.v.UnsafePointer())
	ro.init()

	c := getCodec(rec.t)
	d := compact.NewDeserializer(rec.data)
	c.dec(d, func(a any) {
		v := reflect.ValueOf(a)
		ro.v.Elem().Set(v)
	})
	return ro
}

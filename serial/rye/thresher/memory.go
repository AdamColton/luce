package thresher

import (
	"bytes"
	"crypto/rand"
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
)

var byPtr = lmap.Map[uintptr, *rootObj]{}
var byID = lmap.Map[string, *rootObj]{}

type rootObj struct {
	addr uintptr
	v    reflect.Value
	id   []byte
}

func (ro *rootObj) init() {
	ptr := ro.addr
	byPtr[ptr] = ro
	idStr := string(ro.id)
	byID[idStr] = ro
	// runtime.SetFinalizer(ro.v.Interface(), func(any) {
	// 	delete(byPtr, ptr)
	// 	delete(byID, idStr)
	// })
}

func (ro *rootObj) baseValue() reflect.Value {
	switch ro.v.Kind() {
	case reflect.Slice:
		return ro.v
	}
	return ro.v.Elem()
}

var zeroByte = []byte{0}

func (ro *rootObj) getID() []byte {
	if ro == nil {
		return zeroByte
	}
	return ro.id
}

func newRootObj(ptr uintptr, v reflect.Value) *rootObj {
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

func Get[T any](id []byte) (t T, ok bool) {
	rt := reflector.Type[T]()

	ro := getStoreByID(rt, id)
	if ro == nil || ro.v.Kind() == reflect.Invalid {
		return
	}
	t, ok = ro.v.Interface().(T)
	return
}

func getStoreByID(t reflect.Type, id []byte) *rootObj {
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

	data, found := store[idStr]
	if !found {
		return ro
	}

	str := t.String()
	_ = str
	isPtr := t.Kind() == reflect.Pointer
	if isPtr {
		t = t.Elem()
	}

	ro.v = reflect.New(t)

	c := getBaseCodec(t)
	set := ro.v.Elem()

	switch t.Kind() {
	case reflect.Slice:
		ro.v = ro.v.Elem()
	}
	ro.addr = uintptr(ro.v.UnsafePointer())
	ro.init()

	d := compact.NewDeserializer(data)
	dec := decoders[typeEncoding{
		encID: string(c.encodingID),
		t:     t,
	}]
	v := reflect.ValueOf(dec(d))

	set.Set(v)

	return ro
}

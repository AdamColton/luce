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
	byPtr[ro.addr] = ro
	byID[string(ro.id)] = ro
}

func (ro *rootObj) baseValue() reflect.Value {
	if ro.v.Kind() == reflect.Slice {
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

const idLength = 12

func newRootObj(ptr uintptr, v reflect.Value) *rootObj {
	id := make([]byte, idLength)
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
	d := compact.NewDeserializer(data)
	encID := d.CompactSlice()

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	ro.v = reflect.New(t)
	set := ro.v.Elem()
	if t.Kind() == reflect.Slice {
		ro.v = ro.v.Elem()
	}
	ro.addr = uintptr(ro.v.UnsafePointer())
	ro.init()

	dec := getDecoder(t, encID)
	v := reflect.ValueOf(dec(d))
	set.Set(v)

	return ro
}

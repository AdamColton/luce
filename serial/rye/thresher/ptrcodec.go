package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/serial/rye/compact"
)

var pointerCodec *codec

func initPointerCoded() {
	pointerCodec = &codec{
		enc: func(i any, s compact.Serializer) {
			ro := rootObjByV(reflect.ValueOf(i))
			s.CompactSlice(ro.getID())
		},
		dec: func(d compact.Deserializer) any {
			ro := getStoreByID(d.CompactSlice())
			if ro == nil {
				return nil
			}
			return ro.v.Interface()
		},
		size: func(i any) uint64 {
			ro := rootObjByV(reflect.ValueOf(i))
			return compact.Size(ro.getID())
		},
		roots: func(v reflect.Value) []*rootObj {
			ro := rootObjByV(v)
			if ro != nil {
				return []*rootObj{ro}
			}
			return nil
		},
	}
}

package thresher

import (
	"reflect"
)

type structField struct {
	offset uintptr
	uintPtrOp
	fieldHeader uint64
}

type ptrMarshaller struct {
	op uintPtrOp
	t  reflect.Type
}

type structMarshaller struct {
	byOrder []structField
	byId    map[uint64]structField
}

type marshaller struct {
	op uintPtrOp
	t  reflect.Type
}

type sliceMarshaller struct {
	op        uintPtrOp
	recordLen uintptr
}

type interfaceMarshaller struct {
	t  *Thresher
	rt reflect.Type
}

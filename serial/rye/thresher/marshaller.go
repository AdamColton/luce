package thresher

import (
	"reflect"
)

type structField struct {
	offset uintptr
	UintPtrOp
	fieldHeader uint64
}

type ptrMarshaller struct {
	op UintPtrOp
	t  reflect.Type
}

type structMarshaller struct {
	byOrder []structField
	byId    map[uint64]structField
}

type marshaller struct {
	op UintPtrOp
	t  reflect.Type
}

type sliceMarshaller struct {
	op        UintPtrOp
	recordLen uintptr
}

type interfaceMarshaller struct {
	t  *Thresher
	rt reflect.Type
}

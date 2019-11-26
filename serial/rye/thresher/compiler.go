package thresher

import (
	"reflect"
)

func (t *Thresher) compile(rt reflect.Type) UintPtrOp {
	switch rt.Kind() {
	case reflect.Ptr:
		rt = rt.Elem()
		return ptrMarshaller{
			op: t.compile(rt),
			t:  rt,
		}
	case reflect.Struct:
		return t.compileStruct(rt)
	case reflect.String:
		return uintPtrOpString{}
	case reflect.Int:
		return uintPtrOpInt{}
	case reflect.Int8:
		return uintPtrOpInt8{}
	case reflect.Int16:
		return uintPtrOpInt16C{}
	case reflect.Int32:
		return uintPtrOpInt32C{}
	case reflect.Int64:
		return uintPtrOpInt64C{}
	case reflect.Uint:
		return uintPtrOpUint{}
	case reflect.Uint8:
		return uintPtrOpByte{}
	case reflect.Uint16:
		return uintPtrOpUint16C{}
	case reflect.Uint32:
		return uintPtrOpUint32C{}
	case reflect.Uint64:
		return uintPtrOpUint64C{}
	case reflect.Float32:
		return uintPtrOpFloat32{}
	case reflect.Float64:
		return uintPtrOpFloat64{}
	case reflect.Slice:
		return t.compileSlice(rt.Elem())
	case reflect.Interface:
		return interfaceMarshaller{
			t:  t,
			rt: rt,
		}
	}
	return nil
}

func (t *Thresher) compileStruct(rt reflect.Type) *structMarshaller {
	if t.structMarshallers == nil {
		t.structMarshallers = make(map[reflect.Type]*structMarshaller)
		t.fields = make(map[uint64]field)
	}
	if sm, found := t.structMarshallers[rt]; found {
		return sm
	}
	ln := rt.NumField()
	sm := &structMarshaller{
		byOrder: make([]structField, 0, ln),
		byId:    make(map[uint64]structField, ln),
	}
	t.structMarshallers[rt] = sm
	for i := 0; i < ln; i++ {
		rf := rt.Field(i)
		var skip bool
		f := field{
			name: rf.Name,
			kind: typeID(rf.Type),
		}
		id := f.id()
		sf := structField{
			offset: rf.Offset,
		}
		if skip {
			sf.UintPtrOp = uintPtrOpSkip{}
			sf.fieldHeader = 0
		} else {
			sf.UintPtrOp = t.compile(rf.Type)
			sf.fieldHeader = id
		}
		sm.byOrder = append(sm.byOrder, sf)
		sm.byId[id] = sf
		t.fields[id] = f
	}
	return sm
}

func (t *Thresher) compileSlice(rt reflect.Type) sliceMarshaller {
	return sliceMarshaller{
		recordLen: rt.Size(),
		op:        t.compile(rt),
	}
}

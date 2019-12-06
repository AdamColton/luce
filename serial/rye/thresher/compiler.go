package thresher

import (
	"reflect"
)

type Op struct {
	IsRoot bool
	Base   Compiler
}

func (op Op) Compile(rt reflect.Type) UintPtrOp {
	return op.Base.Compile(rt, op)
}

type Compiler interface {
	Compile(typ reflect.Type, op Op) UintPtrOp
}

type defaultCompiler struct {
	compilers []Compiler
	t         *Thresher
}

func (c defaultCompiler) Compile(rt reflect.Type, op Op) UintPtrOp {
	typeStr := rt.String()
	for _, c := range c.compilers {
		op := c.Compile(rt, op)
		if op != nil {
			return op
		}
	}

	switch rt.Kind() {
	case reflect.Ptr:
		rt = rt.Elem()
		return ptrMarshaller{
			op: op.Compile(rt),
			t:  rt,
		}
	case reflect.Struct:
		return c.compileStruct(rt, op)
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
		return c.compileSlice(rt.Elem(), op)
	case reflect.Interface:
		return interfaceMarshaller{
			t:  c.t,
			rt: rt,
		}
	}
	panic("Not implemented: " + typeStr)
}

func (c defaultCompiler) compileStruct(rt reflect.Type, op Op) *structMarshaller {
	op.IsRoot = false
	str := rt.String()
	_ = str
	if c.t.structMarshallers == nil {
		c.t.structMarshallers = make(map[reflect.Type]*structMarshaller)
		c.t.fields = make(map[uint64]field)
	}
	if sm, found := c.t.structMarshallers[rt]; found {
		return sm
	}
	ln := rt.NumField()
	sm := &structMarshaller{
		byOrder: make([]structField, 0, ln),
		byId:    make(map[uint64]structField, ln),
	}
	c.t.structMarshallers[rt] = sm
	for i := 0; i < ln; i++ {
		rf := rt.Field(i)
		f := field{
			name: rf.Name,
			kind: typeID(rf.Type),
		}
		// if PkgPath is populated, then the field is unexported
		skip := rf.PkgPath != ""
		id := f.id()
		sf := structField{
			offset: rf.Offset,
		}
		if skip {
			sf.UintPtrOp = uintPtrOpSkip{}
			sf.fieldHeader = 0
		} else {
			sf.UintPtrOp = op.Compile(rf.Type)
			sf.fieldHeader = id
		}
		sm.byOrder = append(sm.byOrder, sf)
		sm.byId[id] = sf
		c.t.fields[id] = f
	}
	return sm
}

func (c defaultCompiler) compileSlice(rt reflect.Type, op Op) sliceMarshaller {
	op.IsRoot = false
	return sliceMarshaller{
		recordLen: rt.Size(),
		op:        op.Compile(rt),
	}
}

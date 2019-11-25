package thresher

import (
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/serial/rye"
)

const (
	ifcePtrOffset uintptr = 8
)

type uintPtrOp interface {
	size(u uintptr) int
	marshal(u uintptr, s *rye.Serializer)
	unmarshal(u uintptr, d *rye.Deserializer)
	zero(u uintptr) bool
}

func (p ptrMarshaller) size(u uintptr) int {
	size := 1
	u = *(*uintptr)(unsafe.Pointer(u))
	if u != 0 {
		size += p.op.size(u)
	}
	return size
}

func (p ptrMarshaller) zero(u uintptr) bool {
	return *(*uintptr)(unsafe.Pointer(u)) == 0
}

func (p ptrMarshaller) marshal(u uintptr, s *rye.Serializer) {
	u = *(*uintptr)(unsafe.Pointer(u))
	if u == 0 {
		s.Byte(0)
	} else {
		s.Byte(1)
		p.op.marshal(u, s)
	}
}

func (p ptrMarshaller) unmarshal(u uintptr, d *rye.Deserializer) {
	if d.Byte() == 0 {
		return
	}

	i := reflect.New(p.t).Elem().Interface()
	base := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&i)) + ifcePtrOffset))
	p.op.unmarshal(base, d)
	*(*uintptr)(unsafe.Pointer(u)) = base
}

func (i interfaceMarshaller) zero(u uintptr) bool {
	return *(*uintptr)(unsafe.Pointer(u)) == 0
}

func (i interfaceMarshaller) size(u uintptr) int {
	tid := reflect.NewAt(i.rt, unsafe.Pointer(u)).Elem().Interface().(HasType).TypeID()
	m := i.t.typedIDMarshallers[tid]
	return int(rye.SizeCompactUint64(tid)) + m.op.size(u+ifcePtrOffset)
}

func (i interfaceMarshaller) marshal(u uintptr, s *rye.Serializer) {
	tid := reflect.NewAt(i.rt, unsafe.Pointer(u)).Elem().Interface().(HasType).TypeID()

	m := i.t.typedIDMarshallers[tid]
	s.CompactUint64(tid)
	m.op.marshal(u+ifcePtrOffset, s)
}

func (i interfaceMarshaller) unmarshal(u uintptr, d *rye.Deserializer) {
	tid := d.CompactUint64()
	m := i.t.typedIDMarshallers[tid]
	ifce := reflect.New(m.t).Elem().Interface()
	base := uintptr(unsafe.Pointer(&ifce)) + ifcePtrOffset
	m.op.unmarshal(base, d)

	r := reflect.NewAt(i.rt, unsafe.Pointer(u))
	r.Elem().Set(reflect.ValueOf(ifce))
}

type uintPtrOpByteSlice struct{}

func (uintPtrOpByteSlice) size(u uintptr) int {
	ln := len(*(*[]byte)(unsafe.Pointer(u)))
	return ln + int(rye.SizeCompactUint64(uint64(ln)))
}

func (uintPtrOpByteSlice) zero(u uintptr) bool {
	return len(*(*[]byte)(unsafe.Pointer(u))) == 0
}

func (uintPtrOpByteSlice) marshal(u uintptr, s *rye.Serializer) {
	b := *(*[]byte)(unsafe.Pointer(u))
	s.CompactUint64(uint64(len(b)))
	s.Slice(b)
}
func (uintPtrOpByteSlice) unmarshal(u uintptr, d *rye.Deserializer) {
	ln := int(d.CompactUint64())
	b := (*[]byte)(unsafe.Pointer(u))
	*b = d.Slice(ln)
}

type uintPtrOpString struct{}

func (uintPtrOpString) size(u uintptr) int {
	s := *(*string)(unsafe.Pointer(u))
	ln := len(s)
	return ln + int(rye.SizeCompactUint64(uint64(ln)))
}

func (uintPtrOpString) zero(u uintptr) bool {
	return len(*(*string)(unsafe.Pointer(u))) == 0
}

func (uintPtrOpString) marshal(u uintptr, s *rye.Serializer) {
	str := *(*string)(unsafe.Pointer(u))
	s.CompactUint64(uint64(len(str)))
	s.String(str)
}

func (uintPtrOpString) unmarshal(u uintptr, d *rye.Deserializer) {
	ln := int(d.CompactUint64())
	str := (*string)(unsafe.Pointer(u))
	*str = d.String(ln)
}

func (sm structMarshaller) size(base uintptr) int {
	size := 1
	for _, f := range sm.byOrder {
		if f.fieldHeader == 0 || f.zero(base+f.offset) {
			continue
		}
		size += int(rye.SizeCompactUint64(f.fieldHeader))
		size += f.size(base + f.offset)
	}
	return size
}

func (sm structMarshaller) zero(base uintptr) bool {
	for _, f := range sm.byOrder {
		if f.fieldHeader == 0 {
			continue
		}
		if !f.zero(base + f.offset) {
			return false
		}
	}
	return true
}

func (sm structMarshaller) marshal(base uintptr, s *rye.Serializer) {
	for _, f := range sm.byOrder {
		if f.fieldHeader == 0 || f.zero(base+f.offset) {
			continue
		}
		s.CompactUint64(f.fieldHeader)
		f.marshal(base+f.offset, s)
	}
	s.CompactInt64(0)
}

func (sm structMarshaller) unmarshal(base uintptr, d *rye.Deserializer) {
	for {
		field := d.CompactUint64()
		if field == 0 {
			break
		}
		sf, found := sm.byId[field]
		if !found {
			continue
		}
		sf.unmarshal(base+sf.offset, d)
	}
}

type uintPtrOpSkip struct{}

func (uintPtrOpSkip) size(u uintptr) int {
	return 0
}
func (uintPtrOpSkip) zero(u uintptr) bool {
	return true
}
func (uintPtrOpSkip) marshal(u uintptr, s *rye.Serializer)     {}
func (uintPtrOpSkip) unmarshal(u uintptr, d *rye.Deserializer) {}

func (sm sliceMarshaller) size(base uintptr) int {
	s := *(*[]byte)(unsafe.Pointer(base)) // use []byte, type doesn't actually matter
	ln := uintptr(len(s))
	first := uintptr(unsafe.Pointer(&(s[0])))
	size := int(rye.SizeCompactUint64(uint64(ln)))
	for i := uintptr(0); i < ln; i++ {
		size += sm.op.size(first + i*sm.recordLen)
	}
	return size
}

func (sm sliceMarshaller) zero(base uintptr) bool {
	return len(*(*[]byte)(unsafe.Pointer(base))) == 0
}

func (sm sliceMarshaller) marshal(base uintptr, s *rye.Serializer) {
	l := *(*[]byte)(unsafe.Pointer(base)) // use []byte, type doesn't actually matter
	ln := uintptr(len(l))
	first := uintptr(unsafe.Pointer(&(l[0])))
	s.CompactUint64(uint64(ln))
	for i := uintptr(0); i < ln; i++ {
		sm.op.marshal(first+i*sm.recordLen, s)
	}
}

func (sm sliceMarshaller) unmarshal(u uintptr, d *rye.Deserializer) {
	ln := uintptr(d.CompactUint64())
	s := make([]byte, ln*sm.recordLen)
	first := uintptr(unsafe.Pointer(&(s[0])))
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + 8)) = int(ln)
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + 16)) = int(ln)
	*(*[]byte)(unsafe.Pointer(u)) = s
	for i := uintptr(0); i < ln; i++ {
		sm.op.unmarshal(first+i*sm.recordLen, d)
	}
}

type uintPtrOpFloat32 struct{}

func (uintPtrOpFloat32) size(u uintptr) int {
	return 4
}

func (uintPtrOpFloat32) zero(u uintptr) bool {
	return *(*float32)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpFloat32) marshal(u uintptr, s *rye.Serializer) {
	s.Float32(*(*float32)(unsafe.Pointer(u)))
}
func (uintPtrOpFloat32) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*float32)(unsafe.Pointer(u)) = d.Float32()
}

type uintPtrOpFloat64 struct{}

func (uintPtrOpFloat64) size(u uintptr) int {
	return 8
}

func (uintPtrOpFloat64) zero(u uintptr) bool {
	return *(*float64)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpFloat64) marshal(u uintptr, s *rye.Serializer) {
	s.Float64(*(*float64)(unsafe.Pointer(u)))
}
func (uintPtrOpFloat64) unmarshal(u uintptr, d *rye.Deserializer) {
	*(*float64)(unsafe.Pointer(u)) = d.Float64()
}

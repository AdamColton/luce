package thresher

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/serial/rye"
)

const (
	IfcePtrOffset uintptr = 8
)

type UintPtrOp interface {
	Size(u uintptr) int
	Marshal(u uintptr, s *rye.Serializer)
	Unmarshal(u uintptr, d *rye.Deserializer)
	Zero(u uintptr) bool
}

func (p ptrMarshaller) Size(u uintptr) int {
	size := 1
	u = *(*uintptr)(unsafe.Pointer(u))
	if u != 0 {
		size += p.op.Size(u)
	}
	return size
}

func (p ptrMarshaller) Zero(u uintptr) bool {
	return *(*uintptr)(unsafe.Pointer(u)) == 0
}

func (p ptrMarshaller) Marshal(u uintptr, s *rye.Serializer) {
	u = *(*uintptr)(unsafe.Pointer(u))
	if u == 0 {
		s.Byte(0)
	} else {
		s.Byte(1)
		p.op.Marshal(u, s)
	}
}

func (p ptrMarshaller) Unmarshal(u uintptr, d *rye.Deserializer) {
	if d.Byte() == 0 {
		return
	}

	i := reflect.New(p.t).Elem().Interface()
	base := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&i)) + IfcePtrOffset))
	p.op.Unmarshal(base, d)
	*(*uintptr)(unsafe.Pointer(u)) = base
}

func (i interfaceMarshaller) Zero(u uintptr) bool {
	return *(*uintptr)(unsafe.Pointer(u)) == 0
}

func (i interfaceMarshaller) Size(u uintptr) int {
	tid := reflect.NewAt(i.rt, unsafe.Pointer(u)).Elem().Interface().(HasType).Type()
	idx, _ := i.t.typedIDMarshallersIdx.Get(tid)
	m := i.t.typedIDMarshallers[idx]
	return int(rye.Size(tid)) + m.op.Size(u+IfcePtrOffset)
}

func (i interfaceMarshaller) Marshal(u uintptr, s *rye.Serializer) {
	tid := reflect.NewAt(i.rt, unsafe.Pointer(u)).Elem().Interface().(HasType).Type()
	idx, _ := i.t.typedIDMarshallersIdx.Get(tid)

	m := i.t.typedIDMarshallers[idx]
	s.PrefixSlice(tid)
	m.op.Marshal(u+IfcePtrOffset, s)
}

func (i interfaceMarshaller) Unmarshal(u uintptr, d *rye.Deserializer) {
	tid := d.PrefixSlice()
	idx, _ := i.t.typedIDMarshallersIdx.Get(tid)
	m := i.t.typedIDMarshallers[idx]
	ifce := reflect.New(m.t).Elem().Interface()
	base := uintptr(unsafe.Pointer(&ifce)) + IfcePtrOffset
	m.op.Unmarshal(base, d)

	r := reflect.NewAt(i.rt, unsafe.Pointer(u))
	r.Elem().Set(reflect.ValueOf(ifce))
}

type uintPtrOpByteSlice struct{}

func (uintPtrOpByteSlice) Size(u uintptr) int {
	ln := len(*(*[]byte)(unsafe.Pointer(u)))
	return ln + int(rye.SizeCompactUint64(uint64(ln)))
}

func (uintPtrOpByteSlice) Zero(u uintptr) bool {
	return len(*(*[]byte)(unsafe.Pointer(u))) == 0
}

func (uintPtrOpByteSlice) Marshal(u uintptr, s *rye.Serializer) {
	b := *(*[]byte)(unsafe.Pointer(u))
	s.CompactUint64(uint64(len(b)))
	s.Slice(b)
}
func (uintPtrOpByteSlice) Unmarshal(u uintptr, d *rye.Deserializer) {
	ln := int(d.CompactUint64())
	b := (*[]byte)(unsafe.Pointer(u))
	*b = d.Slice(ln)
}

type uintPtrOpString struct{}

func (uintPtrOpString) Size(u uintptr) int {
	s := *(*string)(unsafe.Pointer(u))
	ln := len(s)
	return ln + int(rye.SizeCompactUint64(uint64(ln)))
}

func (uintPtrOpString) Zero(u uintptr) bool {
	return len(*(*string)(unsafe.Pointer(u))) == 0
}

func (uintPtrOpString) Marshal(u uintptr, s *rye.Serializer) {
	str := *(*string)(unsafe.Pointer(u))
	s.CompactUint64(uint64(len(str)))
	s.String(str)
}

func (uintPtrOpString) Unmarshal(u uintptr, d *rye.Deserializer) {
	ln := int(d.CompactUint64())
	str := (*string)(unsafe.Pointer(u))
	*str = d.String(ln)
}

func (sm structMarshaller) Size(base uintptr) int {
	size := 1
	for _, f := range sm.byOrder {
		if f.fieldHeader == 0 || f.Zero(base+f.offset) {
			continue
		}
		fmt.Println(f.fieldHeader)
		size += int(rye.SizeCompactUint64(f.fieldHeader))
		size += f.Size(base + f.offset)
	}
	return size
}

func (sm structMarshaller) Zero(base uintptr) bool {
	for _, f := range sm.byOrder {
		if f.fieldHeader == 0 {
			continue
		}
		if !f.Zero(base + f.offset) {
			return false
		}
	}
	return true
}

func (sm structMarshaller) Marshal(base uintptr, s *rye.Serializer) {
	for _, f := range sm.byOrder {
		if f.fieldHeader == 0 || f.Zero(base+f.offset) {
			continue
		}
		s.CompactUint64(f.fieldHeader)
		f.Marshal(base+f.offset, s)
	}
	s.CompactInt64(0)
}

func (sm structMarshaller) Unmarshal(base uintptr, d *rye.Deserializer) {
	for {
		field := d.CompactUint64()
		if field == 0 {
			break
		}
		sf, found := sm.byId[field]
		if !found {
			continue
		}
		sf.Unmarshal(base+sf.offset, d)
	}
}

type uintPtrOpSkip struct{}

func (uintPtrOpSkip) Size(u uintptr) int {
	return 0
}
func (uintPtrOpSkip) Zero(u uintptr) bool {
	return true
}
func (uintPtrOpSkip) Marshal(u uintptr, s *rye.Serializer)     {}
func (uintPtrOpSkip) Unmarshal(u uintptr, d *rye.Deserializer) {}

func (sm sliceMarshaller) Size(base uintptr) int {
	s := *(*[]byte)(unsafe.Pointer(base)) // use []byte, type doesn't actually matter
	ln := uintptr(len(s))
	first := uintptr(unsafe.Pointer(&(s[0])))
	size := int(rye.SizeCompactUint64(uint64(ln)))
	for i := uintptr(0); i < ln; i++ {
		size += sm.op.Size(first + i*sm.recordLen)
	}
	return size
}

func (sm sliceMarshaller) Zero(base uintptr) bool {
	return len(*(*[]byte)(unsafe.Pointer(base))) == 0
}

func (sm sliceMarshaller) Marshal(base uintptr, s *rye.Serializer) {
	l := *(*[]byte)(unsafe.Pointer(base)) // use []byte, type doesn't actually matter
	ln := uintptr(len(l))
	first := uintptr(unsafe.Pointer(&(l[0])))
	s.CompactUint64(uint64(ln))
	for i := uintptr(0); i < ln; i++ {
		sm.op.Marshal(first+i*sm.recordLen, s)
	}
}

func (sm sliceMarshaller) Unmarshal(u uintptr, d *rye.Deserializer) {
	ln := uintptr(d.CompactUint64())
	s := make([]byte, ln*sm.recordLen)
	first := uintptr(unsafe.Pointer(&(s[0])))
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + 8)) = int(ln)
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + 16)) = int(ln)
	*(*[]byte)(unsafe.Pointer(u)) = s
	for i := uintptr(0); i < ln; i++ {
		sm.op.Unmarshal(first+i*sm.recordLen, d)
	}
}

type uintPtrOpFloat32 struct{}

func (uintPtrOpFloat32) Size(u uintptr) int {
	return 4
}

func (uintPtrOpFloat32) Zero(u uintptr) bool {
	return *(*float32)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpFloat32) Marshal(u uintptr, s *rye.Serializer) {
	s.Float32(*(*float32)(unsafe.Pointer(u)))
}
func (uintPtrOpFloat32) Unmarshal(u uintptr, d *rye.Deserializer) {
	*(*float32)(unsafe.Pointer(u)) = d.Float32()
}

type uintPtrOpFloat64 struct{}

func (uintPtrOpFloat64) Size(u uintptr) int {
	return 8
}

func (uintPtrOpFloat64) Zero(u uintptr) bool {
	return *(*float64)(unsafe.Pointer(u)) == 0
}
func (uintPtrOpFloat64) Marshal(u uintptr, s *rye.Serializer) {
	s.Float64(*(*float64)(unsafe.Pointer(u)))
}
func (uintPtrOpFloat64) Unmarshal(u uintptr, d *rye.Deserializer) {
	*(*float64)(unsafe.Pointer(u)) = d.Float64()
}

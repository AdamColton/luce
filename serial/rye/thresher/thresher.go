package thresher

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/serial/rye"
)

type Thresher struct {
	typedIDMarshallers []*marshaller
	structMarshallers  map[reflect.Type]*structMarshaller
	fields             map[uint64]field
}

func New() *Thresher {
	return &Thresher{}
}

func (t *Thresher) Unmarshal(data []byte) (HasType, error) {
	d := rye.NewDeserializer(data)
	vt := int(d.CompactUint64())
	if vt > len(t.typedIDMarshallers) {
		return nil, errors.New("Not found")
	}
	m := t.typedIDMarshallers[vt]
	if m == nil {
		return nil, errors.New("Not found")
	}

	r := reflect.New(m.t)
	i := r.Elem().Interface()
	base := uintptr(unsafe.Pointer(&i)) + ifcePtrOffset
	m.op.unmarshal(uintptr(unsafe.Pointer(base)), d)
	return i.(HasType), nil
}

func (t *Thresher) Marshal(v HasType) ([]byte, error) {
	return t.MarshalBuf(v, nil)
}

func (t *Thresher) MarshalBuf(v HasType, in []byte) ([]byte, error) {
	vt := v.TypeID()
	if len(t.typedIDMarshallers) < int(vt) {
		return nil, errors.New("Not found")
	}
	m := t.typedIDMarshallers[vt]
	if m == nil {
		return nil, errors.New("not found")
	}

	base := uintptr(unsafe.Pointer(&v)) + ifcePtrOffset

	s := &rye.Serializer{}
	if in == nil {
		s.Data = make([]byte, m.op.size(base)+int(rye.SizeCompactUint64(vt)))
	} else {
		s.Size = len(in)
		s.Data = in
	}
	s.CompactUint64(vt)
	m.op.marshal(base, s)
	return s.Data, nil
}

func (t *Thresher) Register(vs ...HasType) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, v := range vs {
		err = t.register(v)
		if err != nil {
			return
		}
	}
	return
}

func (t *Thresher) register(v HasType) error {
	vid := v.TypeID()
	if len(t.typedIDMarshallers) <= int(vid) {
		ln := int(vid)
		if ln < 256 {
			ln = 256
		}
		s := make([]*marshaller, ln)
		if t.typedIDMarshallers != nil {
			copy(s, t.typedIDMarshallers)
		}
		t.typedIDMarshallers = s
	}
	if t.typedIDMarshallers[vid] != nil {
		return errors.New("TypeID redefined")
	}
	vt := reflect.TypeOf(v)
	t.typedIDMarshallers[vid] = &marshaller{
		op: t.compile(vt),
		t:  vt,
	}
	return nil
}

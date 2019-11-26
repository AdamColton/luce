package thresher

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"

	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/serial/rye"
)

type Thresher struct {
	typedIDMarshallersIdx byteid.Index
	typedIDMarshallers    []*marshaller
	structMarshallers     map[reflect.Type]*structMarshaller
	fields                map[uint64]field
}

func New() *Thresher {
	return &Thresher{
		typedIDMarshallersIdx: bytebtree.Factory(20),
		typedIDMarshallers:    make([]*marshaller, 20),
	}
}

func (t *Thresher) Unmarshal(data []byte) (HasType, error) {
	d := rye.NewDeserializer(data)
	vt := d.PrefixSlice()
	idx, found := t.typedIDMarshallersIdx.Get(vt)
	if !found {
		return nil, errors.New("Not found")
	}
	m := t.typedIDMarshallers[idx]

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
	vt := v.Type()
	idx, found := t.typedIDMarshallersIdx.Get(vt)
	if !found {
		return nil, errors.New("Not found")
	}
	m := t.typedIDMarshallers[idx]
	if m == nil {
		return nil, errors.New("not found")
	}

	base := uintptr(unsafe.Pointer(&v)) + ifcePtrOffset

	s := &rye.Serializer{}
	if in == nil {
		s.Data = make([]byte, m.op.size(base)+int(rye.Size(vt)))
	} else {
		s.Size = len(in)
		s.Data = in
	}
	s.PrefixSlice(vt)
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

func (th *Thresher) register(v HasType) error {
	vt := v.Type()
	idx, app := th.typedIDMarshallersIdx.Insert(vt)
	if !app && th.typedIDMarshallers[idx] != nil {
		return errors.New("TypeID redefined")
	}
	t := reflect.TypeOf(v)
	m := &marshaller{
		op: th.compile(t),
		t:  t,
	}
	if app {
		th.typedIDMarshallers = append(th.typedIDMarshallers, m)
	} else {
		th.typedIDMarshallers[idx] = m
	}
	return nil
}

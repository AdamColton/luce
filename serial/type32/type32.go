package type32

import (
	"bytes"
	"errors"
	"io"
	"reflect"
)

// TypeIDer32 identifies a type by a uint32. The uint32 size was chosen becuase
// it should allow for plenty of TypeID32 types, but uses little overhead.
type TypeIDer32 interface {
	TypeID32() uint32
}

func uint32ToSlice(u uint32) []byte {
	return []byte{
		byte(u),
		byte(u >> 8),
		byte(u >> 16),
		byte(u >> 24),
	}
}

func sliceToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return uint32(b[0]) + (uint32(b[1]) << 8) + (uint32(b[2]) << 16) + (uint32(b[3]) << 24)
}

// SerializeTypeID32Func is a function signature that, when fulfilled, provides
// a method that fulfills the Serializer signature and handles type id
// prefixing.
type SerializeTypeID32Func func(io.Writer, interface{}) error

// Serialize prepends the TypeID32Type uint32 to a slice then append the
// serialized value. It fulfills the Serialize field on Sender and allows the
// TypeID32Type prefixing strategy to be reused for different serialization
// types.
func (fn SerializeTypeID32Func) Serialize(i interface{}) ([]byte, error) {
	msg, ok := i.(TypeIDer32)
	if !ok {
		return nil, errors.New("Serialize requires interface to be TypeIDer32")
	}

	buf := bytes.NewBuffer(nil)
	buf.Write(uint32ToSlice(msg.TypeID32()))
	err := fn(buf, msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// TypeID32Deserializer fulfills Deserializer and uses the TypeIDer32 prefixing
// strategy.
type TypeID32Deserializer struct {
	types map[uint32]reflect.Type
	fn    func(io.Reader, interface{}) error
}

// DeserializeTypeID32Func function signature, when fulfilled, provides a method
// that creates a TypeID32Deserializer.
type DeserializeTypeID32Func func(io.Reader, interface{}) error

// NewTypeID32Deserializer creates a TypeID32Deserializer from a deserializing
// func.
func (fn DeserializeTypeID32Func) NewTypeID32Deserializer() *TypeID32Deserializer {
	return &TypeID32Deserializer{
		types: make(map[uint32]reflect.Type),
		fn:    fn,
	}
}

// RegisterType with the Deserializer. Fulfills the Deserializer interface.
func (d *TypeID32Deserializer) RegisterType(zeroValue interface{}) error {
	msg, ok := zeroValue.(TypeIDer32)
	if !ok {
		if zeroValue == nil {
			return errors.New("TypeID32Deserializer.Register) cannot register nil interface")
		}
		return errors.New("TypeID32Deserializer.Register) " + reflect.TypeOf(zeroValue).Name() + " does not fulfill TypeID32Type")
	}
	d.types[msg.TypeID32()] = reflect.TypeOf(msg)
	return nil
}

// Deserialize a TypeID32. Fulfills the Deserialize interface.
func (d *TypeID32Deserializer) Deserialize(b []byte) (interface{}, error) {
	if len(b) < 4 {
		return nil, errors.New("TypeID32 too short")
	}

	rt := d.types[sliceToUint32(b)]
	if rt == nil {
		return nil, errors.New("No type registered")
	}
	v := reflect.New(rt)
	i := v.Interface()

	err := d.fn(bytes.NewReader(b[4:]), i)
	if err != nil {
		return nil, err
	}
	return reflect.ValueOf(i).Elem().Interface(), nil
}

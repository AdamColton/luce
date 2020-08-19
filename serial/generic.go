package serial

import (
	"bytes"
	"io"
	"reflect"
)

type wrapPrefixer struct {
	ReflectTypePrefixer
}

func (wp wrapPrefixer) PrefixInterfaceType(i interface{}, b []byte) ([]byte, error) {
	return wp.PrefixReflectType(reflect.TypeOf(i), b)
}

// WrapPrefixer takes a ReflectTypePrefixer and wraps it with logic to add
// PrefixInterfaceType there by fulfilling TypePrefixer. This makes it easy to
// fulfill any type of prefixer by just fulfilling ReflectTypePrefixer.
func WrapPrefixer(pre ReflectTypePrefixer) TypePrefixer {
	p, ok := pre.(TypePrefixer)
	if ok {
		return p
	}
	return wrapPrefixer{pre}
}

// PrefixSerializer combines a Prefixer and a Serializer to create a
// TypeSerializer.
type PrefixSerializer struct {
	InterfaceTypePrefixer
	Serializer
}

// SerializeType fulfills TypeSerializer using the underlying
// InterfaceTypePrefixer to prefix a type and the Serializer to serialize the
// data.
func (ps PrefixSerializer) SerializeType(i interface{}, b []byte) ([]byte, error) {
	var err error
	b, err = ps.PrefixInterfaceType(i, b)
	if err != nil {
		return nil, err
	}

	return ps.Serialize(i, b)
}

// WriterSerializer serializes the provided interface to the Writer.
type WriterSerializer func(interface{}, io.Writer) error

// Serialize the interface to the byte slice.
func (fn WriterSerializer) Serialize(i interface{}, b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	err := fn(i, buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ReaderDeserializer serializes the provided interface to the Reader.
type ReaderDeserializer func(interface{}, io.Reader) error

// Deserialize the interface from the byte slice.
func (fn ReaderDeserializer) Deserialize(i interface{}, b []byte) error {
	buf := bytes.NewBuffer(b)
	return fn(i, buf)
}

// PrefixDeserializer handles both deserializing the type and the data. This
// allows a byte slice to be passed in an the interface returned as opposed to
// the normal requirement of passing in both an interface and a slice.
type PrefixDeserializer struct {
	Detyper
	Deserializer
}

// DeserializeType gets the type from the data, creates an instance and then
// deserializes the data into that instance.
func (ds PrefixDeserializer) DeserializeType(data []byte) (interface{}, error) {
	t, data, err := ds.GetType(data)
	if err != nil {
		return nil, err
	}

	var i interface{}
	var isPtr = t.Kind() == reflect.Ptr
	if isPtr {
		i = reflect.New(t.Elem()).Interface()
	} else {
		i = reflect.New(t).Interface()
	}

	err = ds.Deserialize(i, data)
	if err != nil {
		return nil, err
	}

	if isPtr {
		return i, nil
	}
	return reflect.ValueOf(i).Elem().Interface(), nil
}

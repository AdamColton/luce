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

// WriterSerializer is provided by json and gob in the standard library and may
// be provided by other interfaces.
type WriterSerializer func(io.Writer, interface{}) error

func (fn WriterSerializer) Serialize(i interface{}, b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	err := fn(buf, i)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type ReaderDeserializer func(io.Reader, interface{}) error

func (fn ReaderDeserializer) Deserialize(i interface{}, b []byte) error {
	buf := bytes.NewBuffer(b)
	return fn(buf, i)
}

type PrefixDeserializer struct {
	Detyper
	Deserializer
}

func (ds PrefixDeserializer) DeserializeType(data []byte) (interface{}, error) {
	t, data, err := ds.GetType(data)
	if err != nil {
		return nil, err
	}

	var i interface{}
	if t.Kind() == reflect.Ptr {
		i = reflect.New(t.Elem()).Interface()
	} else {
		i = reflect.New(t).Elem().Interface()
	}

	err = ds.Deserialize(i, data)
	if err != nil {
		return nil, err
	}

	return i, nil
}

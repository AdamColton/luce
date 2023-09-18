package serial

import (
	"reflect"
)

// PrefixDeserializer handles both deserializing the type and the data. This
// allows a byte slice to be passed in and the interface returned, as opposed to
// the normal requirement of passing in both an interface and a slice. It
// fulfills TypeDeserializer.
type PrefixDeserializer struct {
	Detyper
	Deserializer
}

// DeserializeType gets the type from the data, creates an instance and then
// deserializes the data into that instance. If the type is a pointer, the
// pointer is returned, otherwise the value it points to. If either step
// fails, the result is nil and the error is returned.
func (ds PrefixDeserializer) DeserializeType(data []byte) (any, error) {
	t, data, err := ds.GetType(data)
	if err != nil {
		return nil, err
	}

	var i any
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

package serial

import (
	"reflect"
)

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

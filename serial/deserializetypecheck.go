package serial

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

// DetyperDeserializer combines a Detyper and a Deserializer.
type DetyperDeserializer interface {
	Detyper
	Deserializer
}

// DeserializeTypeCheck returns a function that reads the type from the data,
// checks that it is T and deserializes the rest of the data into a new T. If T
// is a pointer type the new pointer is returned, otherwise the deserialized
// value. A type other than T returns a lerr.ErrTypeMismatch.
func DeserializeTypeCheck[T any](ds DetyperDeserializer) func(data []byte) (T, error) {
	typeCheck := lerr.TypeChecker[T]()

	t := reflector.Type[T]()
	var getT func() any
	var isPtr = t.Kind() == reflect.Ptr
	if isPtr {
		getT = func() any { return reflect.New(t.Elem()).Interface() }
	} else {
		getT = func() any { return reflect.New(t).Interface() }
	}

	return func(data []byte) (t T, err error) {
		var gotType reflect.Type
		gotType, data, err = ds.GetType(data)
		if err != nil {
			return
		}
		err = typeCheck(gotType)
		if err != nil {
			return
		}
		ti := getT()
		err = ds.Deserialize(ti, data)
		if err != nil {
			return
		}
		if isPtr {
			t = ti.(T)
		} else {
			t = *ti.(*T)
		}
		return
	}
}

// DeserializeToTypeCheck returns a function that reads the type from the data,
// checks that it is T and deserializes the rest of the data into the provided
// T. A type other than T returns a lerr.ErrTypeMismatch.
func DeserializeToTypeCheck[T any](ds DetyperDeserializer) func(t T, data []byte) error {
	typeCheck := lerr.TypeChecker[T]()

	return func(v T, data []byte) (err error) {
		var gotType reflect.Type
		gotType, data, err = ds.GetType(data)
		if err != nil {
			return
		}
		err = typeCheck(gotType)
		if err != nil {
			return
		}
		err = ds.Deserialize(v, data)
		return
	}
}

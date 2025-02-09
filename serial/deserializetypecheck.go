package serial

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

type DetyperDeserializer interface {
	Detyper
	Deserializer
}

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
		t = ti.(T)
		return
	}
}

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

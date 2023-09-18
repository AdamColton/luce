package serial_test

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

type person struct {
	Name string
	Age  int
}

func mockSerialize(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

func errSerializeFn(i interface{}, w io.Writer) error {
	return errSerialize
}

func mockDeserialize(i interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

type typeMap struct{}

var (
	personPtrType = reflector.Type[*person]()
	personType    = personPtrType.Elem()
	jsonStr       = "{\"Name\":\"Adam\",\"Age\":39}\n"
	testPerson    = person{
		Name: "Adam",
		Age:  39,
	}
	errBadPrefix    = lerr.Str("bad prefix")
	errSerialize    = lerr.Str("serialize error")
	errUnregistered = lerr.Str("unregistered type")
)

func (typeMap) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	if t == personPtrType {
		return append(b, 1), nil
	}
	if t == personType {
		return append(b, 2), nil
	}
	return nil, errUnregistered
}

func (typeMap) GetType(data []byte) (t reflect.Type, rest []byte, err error) {
	if data[0] == 1 {
		return personPtrType, data[1:], nil
	}
	if data[0] == 2 {
		return personType, data[1:], nil
	}
	return personPtrType, nil, errBadPrefix
}

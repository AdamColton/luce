package serial_test

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
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

func TestRoundTrip(t *testing.T) {
	tm := typeMap{}
	s := serial.PrefixSerializer{
		InterfaceTypePrefixer: serial.WrapPrefixer(tm),
		Serializer:            serial.WriterSerializer(mockSerialize),
	}
	d := serial.PrefixDeserializer{
		Detyper:      tm,
		Deserializer: serial.ReaderDeserializer(mockDeserialize),
	}

	b, err := s.SerializeType(&testPerson, nil)
	assert.NoError(t, err)

	got, err := d.DeserializeType(b)
	assert.NoError(t, err)
	assert.Equal(t, &testPerson, got)
}

type mockTypeRegistrar func(zeroValue interface{}) error

func (fn mockTypeRegistrar) RegisterType(zeroValue interface{}) error {
	return fn(zeroValue)
}

func TestRegisterTypes(t *testing.T) {
	seen := make(map[string]bool)
	tr := mockTypeRegistrar(func(zeroValue interface{}) error {
		str := reflect.TypeOf(zeroValue).String()
		seen[str] = true
		return nil
	})

	err := serial.RegisterTypes(tr, (*person)(nil), "", int(0))
	assert.NoError(t, err)

	expected := map[string]bool{
		"string":              true,
		"int":                 true,
		"*serial_test.person": true,
	}
	assert.Equal(t, expected, seen)

	errTest := lerr.Str("test err")
	tr = mockTypeRegistrar(func(zeroValue interface{}) error {
		return errTest
	})
	err = serial.RegisterTypes(tr, 1.1)
	assert.Equal(t, errTest, err)
}

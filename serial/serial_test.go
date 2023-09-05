package serial

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
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

func mockDeserialize(i interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

type typeMap struct{}

var (
	personPtrType = reflector.Type[*person]()
	personType    = personPtrType.Elem()
)

func (typeMap) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	if t == personPtrType {
		return append(b, 1), nil
	}
	if t == personType {
		return append(b, 2), nil
	}
	return nil, lerr.Str("Only supports *person for testing")
}

func (typeMap) GetType(data []byte) (t reflect.Type, rest []byte, err error) {
	if data[0] == 1 {
		return personPtrType, data[1:], nil
	}
	if data[0] == 2 {
		return personType, data[1:], nil
	}
	return personPtrType, nil, lerr.Str("Bad prefix")
}

func TestRoundTrip(t *testing.T) {
	tm := typeMap{}
	s := PrefixSerializer{
		InterfaceTypePrefixer: WrapPrefixer(tm),
		Serializer:            WriterSerializer(mockSerialize),
	}
	d := PrefixDeserializer{
		Detyper:      tm,
		Deserializer: ReaderDeserializer(mockDeserialize),
	}

	p := &person{
		Name: "Adam",
		Age:  35,
	}
	b, err := s.SerializeType(p, nil)
	assert.NoError(t, err)

	got, err := d.DeserializeType(b)
	assert.NoError(t, err)
	assert.Equal(t, p, got)
}

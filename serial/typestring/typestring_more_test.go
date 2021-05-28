package typestring_test

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/typestring"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name string
	Age  int
}

func (person) TypeIDString() string { return "person" }

func encode(i any, w io.Writer) error { return json.NewEncoder(w).Encode(i) }
func decode(i any, r io.Reader) error { return json.NewDecoder(r).Decode(i) }

func TestStringPrefixerNotStringer(t *testing.T) {
	buf, err := typestring.StringPrefixer{}.PrefixInterfaceType(42, []byte("keep"))
	assert.Equal(t, typestring.ErrTypeNotFound, err)
	assert.Nil(t, buf)
}

func TestTypeMapPrefixReflectType(t *testing.T) {
	tm := typestring.NewTypeMap()
	tm.Add(ltype.Int, "Int")

	buf, err := tm.PrefixReflectType(ltype.Int, []byte("x:"))
	assert.NoError(t, err)
	assert.Equal(t, "x:Int ", string(buf))

	_, err = tm.PrefixReflectType(ltype.String, nil)
	assert.Equal(t, typestring.ErrTypeNotFound, err)
}

func TestTypeMapGetTypeErrors(t *testing.T) {
	tm := typestring.NewTypeMap(person{})

	_, _, err := tm.GetType([]byte("person"))
	assert.Equal(t, typestring.ErrMalformed, err)

	_, _, err = tm.GetType([]byte(" person data"))
	assert.Equal(t, typestring.ErrMalformed, err)

	_, _, err = tm.GetType([]byte("unknown data"))
	assert.Equal(t, typestring.ErrNotRegistered, err)
}

func TestTypeMapRoundTrip(t *testing.T) {
	tm := typestring.NewTypeMap(person{})
	want := person{Name: "Adam", Age: 39}

	// The Writer and Reader forms fulfill the same interfaces as the plain
	// forms.
	for _, s := range []serial.PrefixSerializer{
		tm.WriterSerializer(encode),
		tm.Serializer(serial.WriterSerializer(encode)),
	} {
		data, err := s.SerializeType(want, nil)
		assert.NoError(t, err)
		assert.Equal(t, "person ", string(data[:7]))

		for _, d := range []serial.PrefixDeserializer{
			tm.ReaderDeserializer(decode),
			tm.Deserializer(serial.ReaderDeserializer(decode)),
		} {
			got, err := d.DeserializeType(data)
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}
}

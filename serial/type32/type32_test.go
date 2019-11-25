package type32

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/adamcolton/luce/serial"

	"github.com/testify/assert"
)

type person struct {
	Name string
	Age  int
}

func (*person) TypeID32() uint32 {
	return 12345
}

type cannotJson struct {
	Fn func()
}

func (*cannotJson) TypeID32() uint32 {
	return 11111
}

func mockSerialize(w io.Writer, i interface{}) error {
	return json.NewEncoder(w).Encode(i)
}

func mockDeserialize(r io.Reader, i interface{}) error {
	return json.NewDecoder(r).Decode(i)
}

func TestRoundTrip(t *testing.T) {
	s := SerializeTypeID32Func(mockSerialize)
	d := DeserializeTypeID32Func(mockDeserialize).NewTypeID32Deserializer()
	err := d.RegisterType((*person)(nil))
	assert.NoError(t, err)

	p := &person{
		Name: "Adam",
		Age:  35,
	}
	b, err := s.Serialize(p)
	assert.NoError(t, err)

	got, err := d.Deserialize(b)
	assert.NoError(t, err)
	assert.Equal(t, p, got)
}

func TestErrorCases(t *testing.T) {
	d := DeserializeTypeID32Func(mockDeserialize).NewTypeID32Deserializer()
	var t32 TypeIDer32
	err := d.RegisterType(t32)
	assert.Error(t, err)
	err = d.RegisterType(123)
	assert.Error(t, err)
	_, err = d.Deserialize([]byte{3})
	assert.Error(t, err)
	_, err = d.Deserialize([]byte{1, 2, 3, 4})

	err = d.RegisterType((*person)(nil))
	assert.NoError(t, err)
	_, err = d.Deserialize([]byte{57, 48, 0, 0, 1, 2, 3})
	err = d.RegisterType((*person)(nil))

	s := SerializeTypeID32Func(mockSerialize)
	_, err = s.Serialize(123)
	assert.Error(t, err)
	_, err = s.Serialize(&cannotJson{
		Fn: func() { t.Error("should not be invoked") },
	})
	assert.Error(t, err)

	assert.Equal(t, uint32(0), sliceToUint32([]byte{5}))
}

func TestType(t *testing.T) {
	var _ serial.Serializer = SerializeTypeID32Func(mockSerialize)
	var _ serial.Deserializer = DeserializeTypeID32Func(mockDeserialize).NewTypeID32Deserializer()
}

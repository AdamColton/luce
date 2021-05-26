package testutil

import (
	"testing"

	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name string
	Age  int
}

func (*person) TypeID32() uint32 {
	return 12345
}

func Suite(t *testing.T, serialize serial.WriterSerializer, deserialize serial.ReaderDeserializer) {
	tm := type32.NewTypeMap()
	s := tm.Serializer(serial.WriterSerializer(serialize))
	d := tm.Deserializer(serial.ReaderDeserializer(deserialize))

	err := tm.RegisterType((*person)(nil))
	assert.NoError(t, err)

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

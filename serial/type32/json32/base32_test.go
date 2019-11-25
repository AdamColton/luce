package json32

import (
	"testing"

	"github.com/testify/assert"
)

type person struct {
	Name string
	Age  int
}

func (*person) TypeID32() uint32 {
	return 12345
}

func TestRoundTrip(t *testing.T) {
	s := Serializer()
	d := Deserializer()
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

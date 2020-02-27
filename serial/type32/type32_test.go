package type32

import (
	"testing"

	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/testify/assert"
)

type person struct {
	Name string
	Age  int
}

func (*person) TypeID32() uint32 {
	return 12345
}

type strSlice []string

func (strSlice) TypeID32() uint32 {
	return 67890
}

func TestRoundTrip(t *testing.T) {
	tm := NewTypeMap()
	s := tm.WriterSerializer(json.Serialize)
	d := tm.ReaderDeserializer(json.Deserialize)

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

	err = tm.RegisterType((strSlice)(nil))
	assert.NoError(t, err)

	sl := strSlice{"one", "two", "three"}
	b, err = s.SerializeType(sl, nil)
	assert.NoError(t, err)

	got, err = d.DeserializeType(b)
	assert.NoError(t, err)
	assert.Equal(t, sl, got)
}

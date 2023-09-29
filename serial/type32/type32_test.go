package type32_test

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

func TestErrTypeNotFound(t *testing.T) {
	intTp := reflector.Type[int]()
	assert.Equal(t, "Type int was not found", type32.ErrTypeNotFound{intTp}.Error())

}

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

func mockSerialize(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

func mockDeserialize(i interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

func TestRoundTrip(t *testing.T) {
	tm := type32.NewTypeMap()
	s := tm.WriterSerializer(mockSerialize)
	d := tm.ReaderDeserializer(mockDeserialize)

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

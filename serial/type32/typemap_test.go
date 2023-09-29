package type32_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

func TestTypeMap(t *testing.T) {
	tm := type32.NewTypeMap()

	err := tm.RegisterType(nil)
	assert.Equal(t, type32.ErrNilZero, err)

	err = tm.RegisterType("test")
	assert.Equal(t, "TypeID32Deserializer.Register) string does not fulfill TypeID32Type", err.Error())

	p := &person{
		Name: "Adam",
		Age:  39,
	}
	personTp := reflector.Type[*person]()

	err = tm.RegisterType(p)
	assert.NoError(t, err)

	expected := make([]byte, 4)
	rye.Serialize.Uint32(expected, p.TypeID32())

	b, err := tm.PrefixInterfaceType(p, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, b)

	b = b[:0]
	b, err = tm.PrefixReflectType(personTp, b)
	assert.NoError(t, err)
	assert.Equal(t, expected, b)

	expectedRest := []byte{3, 1, 4, 1, 5}
	b = append(b, expectedRest...)
	tp, rest, err := tm.GetType(b)

	assert.Equal(t, personTp, tp)
	assert.Equal(t, expectedRest, rest)
	assert.NoError(t, err)

	_, _, err = tm.GetType(nil)
	assert.Equal(t, type32.ErrTooShort, err)

	_, _, err = tm.GetType([]byte{1, 2, 3, 4})
	assert.Equal(t, type32.ErrNotRegistered, err)

	assert.Equal(t, tm, tm.Serializer(nil).InterfaceTypePrefixer)
	assert.Equal(t, tm, tm.WriterSerializer(nil).InterfaceTypePrefixer)
	assert.Equal(t, tm, tm.Deserializer(nil).Detyper)
	assert.Equal(t, tm, tm.ReaderDeserializer(nil).Detyper)
}

func TestTypeMapRegisterType32s(t *testing.T) {
	tm := type32.NewTypeMap()
	p := &person{}
	var s strSlice
	tm.RegisterType32s(p, s)

	expected := make([]byte, 4)

	rye.Serialize.Uint32(expected, p.TypeID32())
	b, err := tm.PrefixInterfaceType(p, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, b)

	rye.Serialize.Uint32(expected, s.TypeID32())
	b, err = tm.PrefixInterfaceType(s, b[:0])
	assert.NoError(t, err)
	assert.Equal(t, expected, b)
}

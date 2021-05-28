package typestring_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/typestring"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

type mockTypeStringer struct{}

func (mockTypeStringer) TypeIDString() string {
	return "mockTypeStringer"
}

type mockTypeStringer2 struct{}

func (mockTypeStringer2) TypeIDString() string {
	return "AnotherMockTypeStringer"
}

func TestMapPrefixer(t *testing.T) {
	mp := typestring.MapPrefixer{
		ltype.String: "String",
		ltype.Bool:   "Bool",
		ltype.Uint:   "Uint",
		ltype.Int:    "Int",
	}

	buf := []byte("testing:")
	buf, err := mp.PrefixReflectType(ltype.Bool, buf)
	assert.NoError(t, err)
	assert.Equal(t, "testing:Bool ", string(buf))

	buf = []byte("test2:")
	p := mp.Serializer(nil)
	buf, err = p.PrefixInterfaceType("test string", buf)
	assert.NoError(t, err)
	assert.Equal(t, "test2:String ", string(buf))

	_, err = mp.PrefixReflectType(ltype.Float64, buf)
	assert.Equal(t, err, typestring.ErrTypeNotFound)

	mp = nil
	_, err = mp.PrefixReflectType(ltype.String, buf)
	assert.Equal(t, err, typestring.ErrTypeNotFound)
}

func TestStringPrefixer(t *testing.T) {
	sp := typestring.StringPrefixer{}
	mts := mockTypeStringer{}

	buf := []byte("testing:")
	buf, err := sp.PrefixInterfaceType(mts, buf)
	assert.NoError(t, err)
	assert.Equal(t, "testing:mockTypeStringer ", string(buf))

	s := sp.Serializer(nil)
	_, ok := s.InterfaceTypePrefixer.(typestring.StringPrefixer)
	assert.True(t, ok)
}

func TestTypeMap(t *testing.T) {
	tm := typestring.NewTypeMap(mockTypeStringer{})
	tm.Add(ltype.Int, "Int")
	err := tm.RegisterType(mockTypeStringer2{})
	assert.NoError(t, err)

	err = tm.RegisterType(nil)
	assert.Equal(t, typestring.ErrNilZero, err)
	err = tm.RegisterType("")
	expected := lerr.Str("TypeIDStringer.Register) string does not fulfill TypeIDStringer")
	assert.Equal(t, expected, err)

	buf, err := tm.PrefixInterfaceType(mockTypeStringer{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "mockTypeStringer ", string(buf))
	buf = append(buf, []byte("more data")...)

	tp, buf, err := tm.GetType(buf)
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(mockTypeStringer{}), tp)
	assert.Equal(t, "more data", string(buf))
}

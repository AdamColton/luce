package gob_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/serial/wrap/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	gob.Register((*testutil.Person)(nil))
	testutil.SerialFuncsRoundTrip(t, gob.Serialize, gob.Deserialize)
}

func TestErrors(t *testing.T) {
	// A func can't be encoded by gob.
	assert.Error(t, gob.Serialize(func() {}, &bytes.Buffer{}))

	// The data is not gob.
	var i int
	assert.Error(t, gob.Deserialize(&i, bytes.NewBufferString("not gob")))
}

func TestInterfaces(t *testing.T) {
	testutil.SerialInterfacesRoundTrip(t, gob.Serializer{}, gob.Deserializer{})
}

func TestSerializerAppends(t *testing.T) {
	type Name struct {
		Name string
	}
	b, err := gob.Serializer{}.Serialize(Name{Name: "Adam"}, []byte("prefix"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("prefix"), b[:6])

	var got Name
	assert.NoError(t, gob.Deserializer{}.Deserialize(&got, b[6:]))
	assert.Equal(t, "Adam", got.Name)
}

func TestSerializerErrors(t *testing.T) {
	_, err := gob.Serializer{}.Serialize(func() {}, nil)
	assert.Error(t, err)

	var i int
	assert.Error(t, gob.Deserializer{}.Deserialize(&i, []byte("not gob")))
}

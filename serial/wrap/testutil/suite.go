// Package testutil provides testing utilities for fulfilling serial
// interfaeces.
package testutil

import (
	"testing"

	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/stretchr/testify/assert"
)

// Person provides an object to test serialization with.
type Person struct {
	Name string
	Age  int
}

// TypeID32 fulfills TypeIDer32
func (*Person) TypeID32() uint32 {
	return 12345
}

// SerialFuncsRoundTrip takes WriterSerializer and a ReaderDeserializer and
// confirms that they sucessfully perform a round trip.
func SerialFuncsRoundTrip(t *testing.T, serialize serial.WriterSerializer, deserialize serial.ReaderDeserializer) {
	tm := type32.NewTypeMap()
	s := tm.Serializer(serial.WriterSerializer(serialize))
	d := tm.Deserializer(serial.ReaderDeserializer(deserialize))

	err := tm.RegisterType((*Person)(nil))
	assert.NoError(t, err)

	p := &Person{
		Name: "Adam",
		Age:  35,
	}
	b, err := s.SerializeType(p, nil)
	assert.NoError(t, err)

	got, err := d.DeserializeType(b)
	assert.NoError(t, err)
	assert.Equal(t, p, got)
}

// SerialInterfacesRoundTrip takes a Serializer and a Deserializer and
// confirms that they successfully perform a round trip.
func SerialInterfacesRoundTrip(t *testing.T, s serial.Serializer, d serial.Deserializer) {
	p := &Person{
		Name: "Adam",
		Age:  35,
	}

	b := []byte("PREFIX")
	b, err := s.Serialize(p, b)
	assert.NoError(t, err)
	assert.Equal(t, "PREFIX", string(b[:6]))
	b = b[6:]

	got := &Person{}
	err = d.Deserialize(got, b)
	assert.NoError(t, err)
	assert.Equal(t, p, got)
}

// Enc is a func for encoding
type Enc func(any) []byte

// Dec is a func for decoding
type Dec func([]byte, any)

// EncDec does a round trip test using Enc and Dec functions.
func EncDec(t *testing.T, enc Enc, dec Dec) {
	p := &Person{
		Name: "Adam",
		Age:  35,
	}
	got := &Person{}
	dec(enc(p), got)
	assert.Equal(t, p, got)
}

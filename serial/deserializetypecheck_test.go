package serial_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

func TestDeserializeTypeCheck(t *testing.T) {
	pd := serial.PrefixDeserializer{
		Detyper:      typeMap{},
		Deserializer: serial.ReaderDeserializer(mockDeserialize),
	}

	typeChecker := serial.DeserializeTypeCheck[*person](pd)
	data := make([]byte, 1, len(jsonStr)+1)
	data[0] = 1
	data = append(data, []byte(jsonStr)...)

	p, err := typeChecker(data)
	assert.NoError(t, err)
	assert.Equal(t, &testPerson, p)

	data[0] = 2
	p, err = typeChecker(data)
	expectedErr := lerr.ErrTypeMismatch{
		Expected: reflector.Type[*person](),
		Actual:   reflector.Type[person](),
	}
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, p)
}

func TestDeserializeToTypeCheck(t *testing.T) {
	pd := serial.PrefixDeserializer{
		Detyper:      typeMap{},
		Deserializer: serial.ReaderDeserializer(mockDeserialize),
	}
	typeChecker := serial.DeserializeToTypeCheck[*person](pd)

	p := &person{}
	data := make([]byte, 1, len(jsonStr)+1)
	data[0] = 1
	data = append(data, []byte(jsonStr)...)
	err := typeChecker(p, data)
	assert.NoError(t, err)
	assert.Equal(t, &testPerson, p)

	data[0] = 2
	err = typeChecker(p, data)
	expectedErr := lerr.ErrTypeMismatch{
		Expected: reflector.Type[*person](),
		Actual:   reflector.Type[person](),
	}
	assert.Equal(t, expectedErr, err)
}

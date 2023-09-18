package serial_test

import (
	"strings"
	"testing"

	"github.com/adamcolton/luce/serial"
	"github.com/stretchr/testify/assert"
)

func TestPrefixSerializer(t *testing.T) {
	s := serial.PrefixSerializer{
		InterfaceTypePrefixer: serial.WrapPrefixer(typeMap{}),
		Serializer:            serial.WriterSerializer(mockSerialize),
	}

	b, err := s.SerializeType(&testPerson, nil)
	assert.NoError(t, err)
	assert.Equal(t, byte(1), b[0])
	assert.Equal(t, string(b[1:]), jsonStr)

	b, err = s.SerializeType(strings.NewReplacer(), nil)
	assert.Equal(t, errUnregistered, err)
	assert.Nil(t, b)
}

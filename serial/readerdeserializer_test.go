package serial_test

import (
	"testing"

	"github.com/adamcolton/luce/serial"
	"github.com/stretchr/testify/assert"
)

func TestReaderDeserializer(t *testing.T) {
	rd := serial.ReaderDeserializer(mockDeserialize)

	var p person
	err := rd.Deserialize(&p, []byte(jsonStr))
	assert.NoError(t, err)
	assert.Equal(t, testPerson, p)
}

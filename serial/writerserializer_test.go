package serial_test

import (
	"testing"

	"github.com/adamcolton/luce/serial"
	"github.com/stretchr/testify/assert"
)

func TestWriterSerializer(t *testing.T) {
	ws := serial.WriterSerializer(mockSerialize)

	data, err := ws.Serialize(testPerson, nil)
	assert.NoError(t, err)
	assert.Equal(t, jsonStr, string(data))

	ws = serial.WriterSerializer(errSerializeFn)
	data, err = ws.Serialize(testPerson, nil)
	assert.Equal(t, errSerialize, err)
	assert.Nil(t, data)
}

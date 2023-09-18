package serial_test

import (
	"encoding/json"
	"testing"

	"github.com/adamcolton/luce/serial"
	"github.com/stretchr/testify/assert"
)

func TestPrefixDeserializer(t *testing.T) {
	pd := serial.PrefixDeserializer{
		Detyper:      typeMap{},
		Deserializer: serial.ReaderDeserializer(mockDeserialize),
	}

	data := make([]byte, 1, len(jsonStr)+1)
	data[0] = 2
	data = append(data, []byte(jsonStr)...)

	i, err := pd.DeserializeType(data)
	assert.NoError(t, err)
	assert.Equal(t, testPerson, i.(person))

	data[0] = 1
	i, err = pd.DeserializeType(data)
	assert.NoError(t, err)
	assert.Equal(t, &testPerson, i.(*person))

	data[0] = 0
	i, err = pd.DeserializeType(data)
	assert.Equal(t, errBadPrefix, err)
	assert.Nil(t, i)

	data = append(data[:1], []byte("bad json string")...)
	data[0] = 1
	i, err = pd.DeserializeType(data)
	assert.IsType(t, &json.SyntaxError{}, err)
	assert.Nil(t, i)
}

package jsondecoder_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/adamcolton/luce/lhttp/jsondecoder"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name  string
	Age   int
	Admin bool
}

func TestJsonDecoder(t *testing.T) {
	expected := &Person{
		Name:  "Adam",
		Age:   37,
		Admin: true,
	}
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(expected)
	assert.NoError(t, err)

	p := &Person{}
	r := httptest.NewRequest("POST", "/", buf)
	err = jsondecoder.New().Decode(p, r)
	assert.NoError(t, err)
	assert.Equal(t, expected, p)
}

func TestJsonDecoderError(t *testing.T) {
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString("not json"))
	err := jsondecoder.New().Decode(&Person{}, r)
	assert.Error(t, err)
}

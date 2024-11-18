package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestUnsupportedType(t *testing.T) {
	ctx := ljson.NewMarshalContext()
	_, err := ljson.Marshal(make(chan int), ctx)
	assert.Error(t, err)
	_, err = ljson.Stringify(make(chan int), ctx)
	assert.Error(t, err)

	// A MarshalContext without a TypesContext cannot marshal anything.
	_, err = ljson.Marshal("x", &ljson.MarshalContext{})
	assert.Error(t, err)
	_, err = ljson.Stringify("x", &ljson.MarshalContext{})
	assert.Error(t, err)
}

func TestMarshalTwice(t *testing.T) {
	ctx := ljson.NewMarshalContext()
	// The second time, the marshaler for the type is already built.
	for i := 0; i < 2; i++ {
		wn, err := ljson.Marshal("x", ctx)
		assert.NoError(t, err)
		assert.Equal(t, `"x"`, wn.String())
	}
}

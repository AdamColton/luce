package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestSerializeError(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	_, err := ctx.Serialize(make(chan int), nil)
	assert.Error(t, err)
}

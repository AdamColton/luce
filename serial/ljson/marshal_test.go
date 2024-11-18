package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestMarshalString(t *testing.T) {
	ctx := ljson.NewMarshalContext()

	str, err := ljson.Stringify("this is a test", ctx)
	assert.NoError(t, err)
	assert.Equal(t, `"this is a test"`, str)
}

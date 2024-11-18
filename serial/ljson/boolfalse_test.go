package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestMarshalBoolFalse(t *testing.T) {
	str, err := ljson.Stringify(false, ljson.NewMarshalContext(false))
	assert.NoError(t, err)
	assert.Equal(t, "false", str)
}

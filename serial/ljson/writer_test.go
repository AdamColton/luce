package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestStringWriter(t *testing.T) {
	strWriter, err := ljson.MarshalString("this is a test", nil)
	assert.NoError(t, err)
	assert.Equal(t, `"this is a test"`, strWriter.String())
}

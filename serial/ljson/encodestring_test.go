package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestEncodeString(t *testing.T) {
	got := ljson.EncodeString(nil, "this is a test", true)
	assert.Equal(t, `"this is a test"`, string(got))
}

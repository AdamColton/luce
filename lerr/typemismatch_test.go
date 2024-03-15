package lerr_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestTypeMismatch(t *testing.T) {
	var a any = 5
	var b any = "5"
	err := lerr.NewTypeMismatch(a, b)
	assert.Equal(t, `Types do not match: expected "int", got "string"`, err.Error())
	assert.Equal(t, err, lerr.TypeMismatch(a, b))

	b = 7
	err = lerr.NewTypeMismatch(a, b)
	assert.NoError(t, err)
}

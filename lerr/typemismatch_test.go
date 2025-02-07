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

	tc := lerr.TypeChecker[string]()
	err = tc(a)
	assert.Equal(t, `Types do not match: expected "string", got "int"`, err.Error())
	err = tc(b)
	assert.NoError(t, err)

	b = 7
	err = lerr.NewTypeMismatch(a, b)
	assert.NoError(t, err)
}

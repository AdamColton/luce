package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestTypePanic(t *testing.T) {
	isString := filter.IsType(ltype.String)
	assert.Equal(t, ltype.String, isString.Panic("a string"))
	assert.Equal(t, ltype.String, isString.Panic(ltype.String))
	assert.PanicsWithValue(t, "did not get expected type", func() {
		isString.Panic(42)
	})
}

package ltype_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestBuiltinTypes(t *testing.T) {
	assert.Equal(t, reflect.TypeOf(""), ltype.String)
	assert.Equal(t, reflect.TypeOf(true), ltype.Bool)
	assert.Equal(t, reflect.TypeOf(0), ltype.Int)
	assert.Equal(t, reflect.TypeOf(int8(0)), ltype.Int8)
	assert.Equal(t, reflect.TypeOf(int16(0)), ltype.Int16)
	assert.Equal(t, reflect.TypeOf(int32(0)), ltype.Int32)
	assert.Equal(t, reflect.TypeOf(int64(0)), ltype.Int64)
	assert.Equal(t, reflect.TypeOf(uint(0)), ltype.Uint)
	assert.Equal(t, reflect.TypeOf(uint8(0)), ltype.Uint8)
	assert.Equal(t, reflect.TypeOf(uint16(0)), ltype.Uint16)
	assert.Equal(t, reflect.TypeOf(uint32(0)), ltype.Uint32)
	assert.Equal(t, reflect.TypeOf(uint64(0)), ltype.Uint64)
	assert.Equal(t, reflect.TypeOf(float32(0)), ltype.Float32)
	assert.Equal(t, reflect.TypeOf(float64(0)), ltype.Float64)

	// Byte is an alias for uint8.
	assert.Equal(t, ltype.Uint8, ltype.Byte)
	assert.Equal(t, reflect.TypeOf([]byte(nil)), ltype.ByteSlice)

	assert.Equal(t, reflect.Interface, ltype.Err.Kind())
	assert.True(t, reflect.TypeOf(assert.AnError).Implements(ltype.Err))
}

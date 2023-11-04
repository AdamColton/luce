package ltype_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestIsPtrToStruct(t *testing.T) {
	type S struct{ A int }
	f := ltype.IsPtrToStruct.Filter

	assert.True(t, f(reflect.TypeOf(&S{})))
	assert.True(t, f(reflect.TypeOf(&struct{}{})))

	assert.False(t, f(reflect.TypeOf(S{})))
	assert.False(t, f(reflect.TypeOf(1)))
	i := 1
	assert.False(t, f(reflect.TypeOf(&i)))
	ps := &S{}
	assert.False(t, f(reflect.TypeOf(&ps)))
	assert.False(t, f(nil))
}

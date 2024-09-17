package filter_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

func TestFilterAndN(t *testing.T) {
	even := filter.New(func(i int) bool { return i%2 == 0 })
	positive := filter.New(func(i int) bool { return i > 0 })
	small := filter.New(func(i int) bool { return i < 10 })

	f := even.AndN(positive, small)
	assert.True(t, f(4))
	assert.False(t, f(3), "not even")
	assert.False(t, f(-2), "not positive")
	assert.False(t, f(12), "not small")

	alone := even.AndN()
	assert.True(t, alone(2))
	assert.False(t, alone(3))
}

func TestTypeAndN(t *testing.T) {
	ptr := filter.IsKind(reflect.Ptr)
	ptrToInt := filter.IsType(reflector.Type[*int]())

	f := ptr.AndN(ptrToInt)
	assert.True(t, f.Filter(reflector.Type[*int]()))
	assert.False(t, f.Filter(reflector.Type[*string]()))
	assert.False(t, f.Filter(reflector.Type[int]()))

	assert.True(t, ptr.AndN().Filter(reflector.Type[*string]()))
}

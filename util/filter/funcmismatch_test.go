package filter_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestFuncMismatch(t *testing.T) {
	intT, strT := reflect.TypeOf(0), reflect.TypeOf("")
	f := filter.Func([]any{intT}, []any{strT})

	assert.True(t, f.Filter(reflect.TypeOf(func(int) string { return "" })))
	// not a function
	assert.False(t, f.Filter(intT))
	// wrong number of arguments or returns
	assert.False(t, f.Filter(reflect.TypeOf(func(int, int) string { return "" })))
	assert.False(t, f.Filter(reflect.TypeOf(func(int) (string, error) { return "", nil })))
	// wrong argument or return type
	assert.False(t, f.Filter(reflect.TypeOf(func(string) string { return "" })))
	assert.False(t, f.Filter(reflect.TypeOf(func(int) int { return 0 })))

	// A filter.Type and a filter.Filter can be used in place of a reflect.Type.
	g := filter.Func(
		[]any{filter.IsType(intT)},
		[]any{filter.Filter[reflect.Type](func(rt reflect.Type) bool { return rt == strT })},
	)
	assert.True(t, g.Filter(reflect.TypeOf(func(int) string { return "" })))
	assert.False(t, g.Filter(reflect.TypeOf(func(int) int { return 0 })))
}

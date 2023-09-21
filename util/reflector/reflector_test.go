package reflector_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

func ExampleType() {
	t := reflector.Type[string]()
	fmt.Println("t is reflect.Type on", t.String())
	// Output: t is reflect.Type on string
}

func TestToType(t *testing.T) {
	s := "string"
	st := reflect.TypeOf(s)

	assert.Equal(t, st, reflector.ToType(s))
	assert.Equal(t, st, reflector.ToType(st))
}

func ExampleToType() {
	t := reflector.ToType("test")
	fmt.Println("t is reflect.Type on", t.String())

	t2 := reflector.ToType(t)
	fmt.Println("t2 is reflect.Type on", t2.String())
	// Output: t is reflect.Type on string
	// t2 is reflect.Type on string
}

func TestToValue(t *testing.T) {
	s := "foo"
	sv := reflect.ValueOf(s)
	str := reflector.ToValue(s).Type().String()
	assert.Equal(t, "string", str)

	assert.Equal(t, sv.Interface(), reflector.ToValue(s).Interface())
	assert.Equal(t, sv, reflector.ToValue(sv))
}

func ExampleToValue() {
	v := reflector.ToValue("test")
	fmt.Println("v is reflect.Value on", v.Kind())

	v2 := reflector.ToValue(v)
	fmt.Println("v2 is reflect.Value on", v2.Kind())
	// Output: v is reflect.Value on string
	// v2 is reflect.Value on string
}

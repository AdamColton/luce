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

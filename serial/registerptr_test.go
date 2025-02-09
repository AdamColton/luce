package serial_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/stretchr/testify/assert"
)

func TestRegisterPtr(t *testing.T) {
	var seen []reflect.Type
	tr := mockTypeRegistrar(func(zeroValue interface{}) error {
		seen = append(seen, reflect.TypeOf(zeroValue))
		return nil
	})

	assert.NoError(t, serial.RegisterPtr[person](tr))
	assert.Equal(t, []reflect.Type{personPtrType}, seen)

	errTest := lerr.Str("test err")
	tr = mockTypeRegistrar(func(zeroValue interface{}) error {
		return errTest
	})
	assert.Equal(t, errTest, serial.RegisterPtr[person](tr))
}

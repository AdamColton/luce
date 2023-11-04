package linject_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/stretchr/testify/assert"
)

type mockFieldSetterInitilizer struct{}

func (mfsi *mockFieldSetterInitilizer) InitilizeField(fn linject.Func, t reflect.Type) linject.FieldSetter {
	return mfsi
}

func (mfsi *mockFieldSetterInitilizer) Set(args []reflect.Value, field reflect.Value) (func(), error) {
	field.Set(reflect.ValueOf("set in Set"))
	return nil, nil
}

func foofunc(s string, data *struct {
	TestField string
}) (string, string) {
	return s, data.TestField
}

func TestFieldSetter(t *testing.T) {
	mfsi := &mockFieldSetterInitilizer{}
	fi := linject.NewFieldInitilizer(mfsi, "TestField")
	fi.FieldType = filter.IsKind(reflect.String)

	m := linject.FuncInitilizers{fi}
	sfn := m.Apply(foofunc).Interface().(func(string) (string, string))
	a, b := sfn("hello")
	assert.Equal(t, "hello", a)
	assert.Equal(t, "set in Set", b)
}

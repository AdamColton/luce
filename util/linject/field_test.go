package linject_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/stretchr/testify/assert"
)

type mockFieldSetterInitilizer struct{}

func (mfsi *mockFieldSetterInitilizer) InitilizeField(fn linject.FuncType, t reflect.Type) linject.FieldInjector {
	return mfsi
}

func (mfsi *mockFieldSetterInitilizer) InjectField(args []reflect.Value, field reflect.Value) (func([]reflect.Value), error) {
	field.Set(reflect.ValueOf("set in InjectField"))
	return nil, nil
}

func foofunc(s string, data *struct {
	TestField string
}) (string, string) {
	return s, data.TestField
}

func TestFieldSetter(t *testing.T) {
	mfsi := &mockFieldSetterInitilizer{}
	fi := linject.NewField(mfsi, "TestField")
	fi.FieldType = filter.IsKind(reflect.String)

	m := linject.Initilizers{fi}
	sfn := m.Apply(foofunc).Interface().(func(string) (string, string))
	a, b := sfn("hello")
	assert.Equal(t, "hello", a)
	assert.Equal(t, "set in InjectField", b)
}

func TestNewFieldSetter(t *testing.T) {
	fn := func() (any, func([]reflect.Value), error) {
		return "set in field setter func", nil, nil
	}
	strType := filter.IsKind(reflect.String)
	anyType := filter.Type{}
	fi := linject.NewFieldSetter(fn, "TestField", anyType, strType)

	m := linject.Initilizers{fi}
	sfn := m.Apply(foofunc).Interface().(func(string) (string, string))
	a, b := sfn("hello")
	assert.Equal(t, "hello", a)
	assert.Equal(t, "set in field setter func", b)
}

package linject_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/stretchr/testify/assert"
)

type nilFieldInitilizer struct{}

func (nilFieldInitilizer) InitilizeField(linject.FuncType, reflect.Type) linject.FieldInjector {
	return nil
}

func TestFieldInitilizeSkips(t *testing.T) {
	mfsi := &mockFieldSetterInitilizer{}
	missing := linject.NewField(mfsi, "Missing")

	wrongField := linject.NewField(mfsi, "TestField")
	wrongField.FieldType = filter.IsKind(reflect.Int)

	wrongFunc := linject.NewField(mfsi, "TestField")
	wrongFunc.FuncType = filter.IsKind(reflect.Int)

	noSetter := linject.NewField(nilFieldInitilizer{}, "TestField")

	for n, fi := range map[string]linject.Field{
		"missing field":  missing,
		"wrong field":    wrongField,
		"wrong function": wrongFunc,
		"no injector":    noSetter,
	} {
		t.Run(n, func(t *testing.T) {
			sfn := linject.Initilizers{fi}.Apply(foofunc).Interface().(func(string) (string, string))
			a, b := sfn("hello")
			assert.Equal(t, "hello", a)
			assert.Empty(t, b, "the field was not set")
		})
	}
}

func TestFieldInvalidName(t *testing.T) {
	fi := linject.NewField(&mockFieldSetterInitilizer{}, "lower")
	assert.Panics(t, func() {
		linject.Initilizers{fi}.Apply(foofunc)
	})
}

func TestNewFieldInjectorArguments(t *testing.T) {
	args := []reflect.Value{reflect.ValueOf("hello")}
	field := func() reflect.Value { return reflect.New(reflect.TypeOf("")).Elem() }

	all := linject.NewFieldInjector(func(args []reflect.Value) (string, func([]reflect.Value), error) {
		return "all: " + args[0].String(), nil, nil
	})
	f := field()
	cb, err := all.InjectField(args, f)
	assert.NoError(t, err)
	assert.Nil(t, cb)
	assert.Equal(t, "all: hello", f.String(), "a single []reflect.Value argument gets them all")

	leading := linject.NewFieldInjector(func(s string) (string, func([]reflect.Value), error) {
		return "leading: " + s, nil, nil
	})
	f = field()
	_, err = leading.InjectField(args, f)
	assert.NoError(t, err)
	assert.Equal(t, "leading: hello", f.String(), "otherwise the leading arguments are matched")
}

func TestNewFieldInjectorError(t *testing.T) {
	boom := errors.New("boom")
	called := false
	inj := linject.NewFieldInjector(func() (string, func([]reflect.Value), error) {
		return "x", func([]reflect.Value) { called = true }, boom
	})
	f := reflect.New(reflect.TypeOf("")).Elem()
	cb, err := inj.InjectField(nil, f)
	assert.Equal(t, boom, err)
	assert.Equal(t, "x", f.String())
	cb(nil)
	assert.True(t, called, "the callback is returned with the error")
}

func TestNewFieldInjectorBadFunc(t *testing.T) {
	assert.Panics(t, func() {
		linject.NewFieldInjector(func() {})
	})
}

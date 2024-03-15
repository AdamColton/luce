package linject_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/linject"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
	Age  int
}

func FooFunc(t *testing.T, data *struct {
	*Person
}) string {
	assert.Equal(t, data.Person.Name, "Adam")
	return "OutString"
}

type PersonInitilizer struct{}

func (pi PersonInitilizer) Initilize(ft linject.FuncType) linject.Injector {
	return personInjector{
		n: ft.Fn().NumIn() - 1,
	}
}

type personInjector struct {
	n int
}

func (pi personInjector) Inject(args []reflect.Value) (callback func([]reflect.Value), err error) {
	d := args[pi.n].Elem()
	p := &Person{
		Name: "Adam",
		Age:  39,
	}
	d.FieldByName("Person").Set(reflect.ValueOf(p))
	cb := func(rets []reflect.Value) {
		s := rets[0].Interface().(string)
		s += " Callback"
		rets[0] = reflect.ValueOf(s)
	}
	return cb, nil
}

func TestFunc(t *testing.T) {
	fi := linject.Initilizers{
		PersonInitilizer{},
	}
	fn := fi.Apply(FooFunc).Interface().(func(*testing.T) string)
	got := fn(t)
	assert.Equal(t, "OutString Callback", got)
}

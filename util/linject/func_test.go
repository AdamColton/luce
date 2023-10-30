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
}) {
	assert.Equal(t, data.Person.Name, "Adam")
}

type PersonInitilizer struct{}

func (pi PersonInitilizer) Initilize(ft linject.Func) linject.DataInserter {
	return personInserter{
		n: ft.Fn().NumIn() - 1,
	}
}

type personInserter struct {
	n int
}

func (pi personInserter) Insert(args []reflect.Value) (callback func(), err error) {
	d := args[pi.n].Elem()
	p := &Person{
		Name: "Adam",
		Age:  39,
	}
	d.FieldByName("Person").Set(reflect.ValueOf(p))
	return nil, nil
}

func TestFunc(t *testing.T) {
	fi := linject.FuncInitilizers{
		PersonInitilizer{},
	}
	fn := fi.Apply(FooFunc).Interface().(func(*testing.T))
	fn(t)
}

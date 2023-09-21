package reflector_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

type foo struct{}

func (foo) A(a, b int)                   {}
func (foo) Hello() string                { return "Hello" }
func (foo) Goodbye() string              { return "Goodbye" }
func (foo) AddExPt(name string) string   { return name + "!" }
func (foo) TwoArgs(first, second string) {}

func TestMethod(t *testing.T) {
	var f foo

	m := reflector.MethodOn(f, "DoesNotExist")
	assert.Nil(t, m)

	m = reflector.MethodOn(f, "Hello")
	assert.Equal(t, "<func() string Value>", m.Func.String())
	assert.Equal(t, "<func(reflector_test.foo) string Value>", m.Method.Func.String())
}

func ExampleMethod_AssignTo() {
	var f foo
	// foo has methods:
	// * A(a, b int)
	// * Hello() string
	// * Goodbye() string
	// * AddExPt(name string) string
	// * TwoArgs(first, second string)

	var fn func(string) string
	for _, m := range reflector.MethodsOn(f) {
		if m.AssignTo(&fn) {
			break
		}
	}

	fmt.Println(fn("goodbye"))
	// Output: goodbye!
}

func ExampleMethods_Funcs() {
	var f foo
	// foo has methods:
	// * A(a, b int)
	// * Hello() string
	// * Goodbye() string
	// * AddExPt(name string) string
	// * TwoArgs(first, second string)

	var seeking = reflector.Type[func(foo) string]()
	var ms reflector.Methods
	for _, m := range reflector.MethodsOn(f) {
		// note: comparing to Method.Type not Func.Type().
		if m.Type == seeking {
			ms = append(ms, m)
		}
	}
	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Name > ms[j].Name
	})

	for _, f := range ms.Funcs() {
		fmt.Println(f.(func() string)())
	}

	// Output:
	// Hello
	// Goodbye
}

func ExampleMethodsOn() {
	var f foo
	// foo has methods:
	// * A(a, b int)
	// * Hello() string
	// * Goodbye() string
	// * AddExPt(name string) string
	// * TwoArgs(first, second string)

	m := reflector.MethodOn(f, "Hello")
	fmt.Println(m.Func.Call(nil)[0])
	// Output: Hello
}

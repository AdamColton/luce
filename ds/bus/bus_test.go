package bus_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/util/handler"
	"github.com/stretchr/testify/assert"
)

type mockListener struct {
	handler.Switcher
	running    bool
	stop       chan bool
	in         <-chan any
	registered []string
}

func (ml *mockListener) Run() {
	ml.running = true
	<-ml.stop
	ml.running = false
}

func (ml *mockListener) SetIn(in <-chan any) {
	ml.in = in
}

func (ml *mockListener) RegisterHandlers(handler ...any) error {
	return nil
}

func (ml *mockListener) RegisterType(zeroValue any) error {
	ml.registered = append(ml.registered, reflect.TypeOf(zeroValue).String())
	return nil
}

func (ml *mockListener) SetErrorHandler(errHandler any) error {
	return nil
}

type handlerObj struct {
}

func (ho handlerObj) StringHandler(s string) {
}

func (ho handlerObj) IntHandler(i int) {
}

func (ho handlerObj) Float64Handler(f float64) {
}

func (ho handlerObj) NotAHandler(i int, s string) {
}

func TestListenerMethodsRegistrar(t *testing.T) {
	ml := &mockListener{
		Switcher: handler.NewSwitch(10),
		stop:     make(chan bool),
	}
	ho := handlerObj{}
	err := bus.DefaultRegistrar.Register(ml, ho)
	assert.NoError(t, err)
	sort.Strings(ml.registered)

	expected := []string{"float64", "int", "string"}
	assert.Equal(t, expected, ml.registered)
}

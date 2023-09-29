package bus_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/stretchr/testify/assert"
)

type mockListener struct {
	handler.Switcher
	running     bool
	stop        chan bool
	in          <-chan any
	registered  []string
	registerErr error
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
	return ml.registerErr
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

// mixedObj has a method that is not a valid handler.
type mixedObj struct{}

func (mixedObj) GoodHandler(s string)        {}
func (mixedObj) BadHandler(i int) (int, int) { return i, i }

func TestRegisterInvalidMethod(t *testing.T) {
	ml := &mockListener{Switcher: handler.NewSwitch(10)}
	err := bus.DefaultRegistrar.Register(ml, mixedObj{})
	assert.Error(t, err)

	// The valid method is registered with both the switch and the receiver.
	assert.Equal(t, []string{"string"}, ml.registered)
	_, err = ml.Handle("hello")
	assert.NoError(t, err)
	_, err = ml.Handle(5)
	assert.Equal(t, handler.ErrNoHandler, err)
}

func TestRegisterTypeError(t *testing.T) {
	ml := &mockListener{
		Switcher:    handler.NewSwitch(10),
		registerErr: lerr.Str("no room"),
	}
	err := bus.DefaultRegistrar.Register(ml, handlerObj{})
	assert.Equal(t, lerr.Str("no room"), err)

	// It stops at the first type that can't be registered.
	assert.Len(t, ml.registered, 1)
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

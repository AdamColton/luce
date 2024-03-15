package linject_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/linject"
	"github.com/stretchr/testify/assert"
)

func TestApplyNotInjectable(t *testing.T) {
	fi := linject.Initilizers{PersonInitilizer{}}
	assert.Nil(t, fi.Apply(func() {}), "no arguments")
	assert.Nil(t, fi.Apply(func(i int) {}), "the last argument is not a pointer")
	assert.Nil(t, fi.Apply(func(p *int) {}), "a pointer, but not to a struct")
}

type typeRecorder struct {
	fn, target reflect.Type
}

func (r *typeRecorder) Initilize(ft linject.FuncType) linject.Injector {
	r.fn, r.target = ft.Fn(), ft.Target()
	return nil
}

func TestFuncType(t *testing.T) {
	r := &typeRecorder{}
	linject.Initilizers{r}.Apply(FooFunc)
	assert.Equal(t, reflect.TypeOf(FooFunc), r.fn)
	assert.Equal(t, reflect.TypeOf(struct{ *Person }{}), r.target)
}

type failingInjector struct {
	err error
	cb  func([]reflect.Value)
}

func (f failingInjector) Initilize(linject.FuncType) linject.Injector { return f }

func (f failingInjector) Inject([]reflect.Value) (func([]reflect.Value), error) {
	return f.cb, f.err
}

func TestInjectors(t *testing.T) {
	boom := errors.New("boom")
	cb := func([]reflect.Value) {}
	injs := linject.Injectors{
		failingInjector{err: boom},
		failingInjector{cb: cb},
		failingInjector{},
	}

	cbs, err := injs.Inject(nil)
	assert.Len(t, cbs, 1, "only the injector that returned a callback adds one")
	assert.ErrorContains(t, err, "boom")

	cbs, err = linject.Injectors{failingInjector{cb: cb}}.Inject(nil)
	assert.Len(t, cbs, 1)
	assert.NoError(t, err)
}

package upgrade_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

type Fooer interface {
	Foo() int
}

type Barer interface {
	Bar() string
}

type Core string

func (c Core) Bar() string {
	return string(c)
}

func (c Core) Foo() int {
	return len(c)
}

type FooWrapper struct {
	Fooer
	Offset int
}

func (f FooWrapper) Foo() int {
	return f.Fooer.Foo() + f.Offset
}

func (f FooWrapper) Upgrade(t reflect.Type) interface{} {
	return upgrade.Wrapped(f.Fooer, t)
}

func TestWrapped(t *testing.T) {
	c := Core("core")
	var f Fooer = FooWrapper{
		Fooer:  c,
		Offset: 2,
	}

	_, ok := f.(Barer)
	assert.False(t, ok)

	var i interface{} = c
	BarerType := reflector.Type[Barer]()
	assert.True(t, reflect.ValueOf(i).Type().Implements(BarerType))

	b := upgrade.Wrapped(f, BarerType).(Barer)
	assert.NotNil(t, b)

	stringerType := reflector.Type[fmt.Stringer]()
	s := upgrade.Wrapped(f, stringerType)
	assert.Nil(t, s)
}

func TestHelper(t *testing.T) {
	c := Core("core")
	var f Fooer = FooWrapper{
		Fooer:  c,
		Offset: 2,
	}

	_, ok := f.(Barer)
	assert.False(t, ok)

	var b Barer
	assert.True(t, upgrade.Upgrade(f, &b))
	assert.NotNil(t, b)
	assert.Equal(t, c.Bar(), b.Bar())

	var s fmt.Stringer
	assert.False(t, upgrade.Upgrade(f, &s))
	assert.Nil(t, s)
}

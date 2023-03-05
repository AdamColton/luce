package upgrade_test

import (
	"fmt"
	"testing"

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

func (f FooWrapper) Wrapped() any {
	return f.Fooer
}

func TestTo(t *testing.T) {
	c := Core("core")
	var f Fooer = FooWrapper{
		Fooer:  c,
		Offset: 2,
	}

	_, ok := f.(Barer)
	assert.False(t, ok)

	b, ok := upgrade.To[Barer](f)
	assert.True(t, ok)
	assert.NotNil(t, b)
	assert.Equal(t, c.Bar(), b.Bar())

	s, ok := upgrade.To[fmt.Stringer](f)
	assert.False(t, ok)
	assert.Nil(t, s)
}

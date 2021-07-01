package reflector

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

type foo struct{}

func (foo) A(a, b int)                   {}
func (foo) Hello()                       {}
func (foo) Goodbye(name string)          {}
func (foo) TwoArgs(first, second string) {}

func TestGetByName(t *testing.T) {
	var f foo
	mf := MethodName(filter.EQ.String("Hello"))
	m := mf.One(f)
	assert.NotNil(t, m)
	var fn func()
	ok := m.SetTo(&fn)
	assert.True(t, ok)
}

func TestGetMany(t *testing.T) {
	var f foo
	mf := ArgCount(filter.LT.Int(2)).MethodFilter()
	ms := mf.On(f)
	assert.Len(t, ms, 2)

	var fs []interface{} = ms.Funcs()
	assert.Len(t, fs, 2)
}

func TestMethodOutOfBounds(t *testing.T) {
	var f foo
	m := NewMethod(reflect.ValueOf(f), 10)
	assert.Nil(t, m)
}

func TestMethodOr(t *testing.T) {
	m1 := FuncOf((func())(nil)).MethodFilter()
	m2 := FuncOf((func(string))(nil)).MethodFilter()
	mf := m1.Or(m2)

	var f foo
	ms := mf.On(f)
	assert.Len(t, ms, 2)
}

func TestMethodAnd(t *testing.T) {
	m1 := ArgCount(filter.EQ.Int(0)).MethodFilter()
	m2 := MethodName(filter.Prefix("He"))
	mf := m1.And(m2)

	var f foo
	ms := mf.On(f)
	assert.Len(t, ms, 1)
	m := ms[0]
	assert.Equal(t, "Hello", m.Name)
}

func TestMethodNot(t *testing.T) {
	mf := MethodName(filter.EQ.String("Hello")).Not()

	var f foo
	ms := mf.On(f)
	assert.Len(t, ms, 3)
	for _, m := range ms {
		assert.NotEqual(t, "Hello", m.Name)
	}
}

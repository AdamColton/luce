package reflector_test

import (
	"testing"

	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name string
	Age  int
}

func TestParser(t *testing.T) {
	p := reflector.Parser[string]{
		reflector.Type[*string]():  reflector.Parsers.String,
		reflector.Type[*float64](): reflector.Parsers.Float64,
		reflector.Type[*int]():     reflector.Parsers.Int,
	}

	var s string
	err := p.Parse(&s, "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", s)

	var i int
	err = p.Parse(&i, "123")
	assert.NoError(t, err)
	assert.Equal(t, 123, i)

	var f float64
	err = p.Parse(&f, "3.14")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, f)

	psn := &person{}
	err = p.ParseFieldName(psn, "Name", "Adam")
	assert.NoError(t, err)
	assert.Equal(t, "Adam", psn.Name)

	err = p.ParseFieldName(psn, "Age", "39")
	assert.NoError(t, err)
	assert.Equal(t, 39, psn.Age)

	// var set string
	// from := "test"
	// sv := reflect.ValueOf(&set)
	// fv := reflect.ValueOf(from)
	// sv.Elem().Set(fv)
	// assert.Equal(t, from, set)
}

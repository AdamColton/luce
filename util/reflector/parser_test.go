package reflector_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/parsers"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	uintPtr := reflector.Type[*uint]()
	strPtr := reflector.Type[*string]()

	p := reflector.Parser[string]{}
	reflector.ParserAdd(p, parsers.String)
	reflector.ParserAdd(p, parsers.Float64)
	reflector.ParserAdd(p, parsers.Int)

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

	err = p.Parse(s, "test")
	assert.Equal(t, reflector.ErrExpectedPtr, err)

	var u uint
	err = p.Parse(&u, "test")
	assert.Equal(t, reflector.ErrParserNotFound{uintPtr}, err)
	assert.Equal(t, "parser not found: *uint", err.Error())

	pf := reflector.ParserFunc[string, *string](parsers.String)
	err = pf.Parser(reflect.ValueOf(&u), "should error")
	expectErr := lerr.TypeMismatch(strPtr, uintPtr)
	assert.Equal(t, expectErr, err)
}

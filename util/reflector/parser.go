package reflector

import (
	"reflect"
	"strings"

	"github.com/adamcolton/luce/lerr"
)

const (
	// ErrParserSet is returned when reflect.Value.Set fails
	ErrParserSet = lerr.Str("could not set")
	// ErrExpectedPtr is returned if the type is not a pointer
	ErrExpectedPtr = lerr.Str("expected ptr")
)

// ErrParserNotFound is returned when a parser does not contain a given Type.
type ErrParserNotFound struct {
	Type reflect.Type
}

// Error fullfils the error interface.
func (err ErrParserNotFound) Error() string {
	return strings.Join([]string{"parser not found: ", err.Type.String()}, "")
}

// Parser is used to parse one data type into many. The most common types for
// T will be string or []byte.
type Parser[T any] map[reflect.Type]func(reflect.Value, T) error

// Parse 't' into interface 'i' using the Parser.
func (p Parser[T]) Parse(i any, t T) error {
	v := reflect.ValueOf(i)
	return p.ParseValue(v, t)
}

// Parse 't' into 'v' using the Parser.
func (p Parser[T]) ParseValue(v reflect.Value, t T) error {
	k := v.Kind()
	if k != reflect.Ptr {
		return ErrExpectedPtr
	}
	tp := v.Type()
	fn, found := p[tp]
	if !found {
		return ErrParserNotFound{tp}
	}
	return fn(v, t)
}

type ParserFunc[In, Out any] func(Out, In) error

func (pf ParserFunc[In, Out]) Parser(v reflect.Value, in In) (err error) {
	out, ok := v.Interface().(Out)
	if !ok {
		return lerr.TypeMismatch(Type[Out](), v.Type())
	}
	return pf(out, in)
}

func ParserAdd[In, Out any](p Parser[In], fn func(Out, In) error) {
	pf := ParserFunc[In, Out](fn)
	p[Type[Out]()] = pf.Parser
}

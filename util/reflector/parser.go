package reflector

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/adamcolton/luce/lerr"
)

const (
	ErrParserNotFound = lerr.Str("parser not found")
	ErrParserSet      = lerr.Str("could not set")
	ErrExpectedPtr    = lerr.Str("expected ptr")
	ErrExpectedStruct = lerr.Str("expected struct")
)

type Parser[T any] map[reflect.Type]func(reflect.Value, T) error

func (p Parser[T]) Parse(i any, t T) error {
	v := reflect.ValueOf(i)
	return p.ParseValue(v, t)
}

func (p Parser[T]) ParseValue(v reflect.Value, t T) error {
	k := v.Kind()
	if k != reflect.Ptr {
		return ErrExpectedPtr
	}
	fn, found := p[v.Type()]
	if !found {
		return ErrParserNotFound
	}
	return fn(v, t)
}

func (p Parser[T]) ParseFieldName(i any, name string, t T) error {
	v := reflect.ValueOf(i)
	return p.ParseValueFieldName(v, name, t)
}

func (p Parser[T]) ParseValueFieldName(v reflect.Value, name string, t T) error {
	if v.Kind() != reflect.Ptr {
		return ErrExpectedPtr
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}
	sf, found := v.Type().FieldByName(name)
	if !found {
		return fmt.Errorf("field '%s' not found", name)
	}
	fv := v.FieldByName(name)
	f := reflect.NewAt(sf.Type, unsafe.Pointer(fv.UnsafeAddr()))
	return p.ParseValue(f, t)
}

type parsers struct{}

var Parsers parsers

func (parsers) String(v reflect.Value, s string) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrParserSet
		}
	}()
	v.Elem().Set(reflect.ValueOf(s))
	return
}

func (parsers) Float64(v reflect.Value, s string) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrParserSet
		}
	}()
	var f float64
	f, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return
	}
	v.Elem().Set(reflect.ValueOf(f))
	return
}

func (parsers) Int(v reflect.Value, s string) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrParserSet
		}
	}()
	var i int
	i, err = strconv.Atoi(s)
	if err != nil {
		return
	}
	v.Elem().Set(reflect.ValueOf(i))
	return
}

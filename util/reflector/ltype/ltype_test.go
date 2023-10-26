package ltype_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestCheckStructPtr(t *testing.T) {
	s := struct {
		Name string
		Age  int
	}{
		Name: "Adam",
		Age:  40,
	}

	tp := ltype.CheckStructPtr(reflect.TypeOf(&s))
	assert.Equal(t, reflect.TypeOf(s), tp)

	tp = ltype.CheckStructPtr(reflect.TypeOf(s))
	assert.Nil(t, tp)
	tp = ltype.CheckStructPtr(reflect.TypeOf(1234))
	assert.Nil(t, tp)
	tp = ltype.CheckStructPtr(reflect.TypeOf(&(s.Name)))
	assert.Nil(t, tp)
}

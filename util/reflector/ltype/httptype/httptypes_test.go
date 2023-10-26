package httptype_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/reflector/ltype/httptype"
	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	assert.Equal(t, reflect.TypeOf((*http.ResponseWriter)(nil)).Elem(), httptype.ResponseWriter)
	assert.Equal(t, reflect.TypeOf(&http.Request{}), httptype.Request)
	assert.Equal(t, reflect.TypeOf(http.HandlerFunc(nil)), httptype.HandlerFunc)
	assert.Equal(t, reflect.TypeOf((*http.Handler)(nil)).Elem(), httptype.Handler)
	assert.Equal(t, reflect.TypeOf(http.Header{}), httptype.Header)
	assert.Equal(t, reflect.TypeOf(&http.Client{}), httptype.Client)
	assert.Equal(t, reflect.TypeOf(&http.Server{}), httptype.Server)
}

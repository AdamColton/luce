package midware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

type mockCallback struct {
	calledBack bool
}

func (m *mockCallback) Initilize(t reflect.Type) DataInserter {
	return m
}

func (m *mockCallback) Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) (func(), error) {
	return m.Callback, nil
}

func (m *mockCallback) Callback() {
	m.calledBack = true
}

type testFieldType struct {
	A string
}

func TestMidware(t *testing.T) {
	mrd := &mockRequestDecoder{
		str: "magic decoder test",
	}
	d := NewDecoder(mrd, "TestField")
	c := &mockCallback{}
	m := New(d, c)
	didRun := false
	fn := m.Handle(func(w http.ResponseWriter, r *http.Request, data *struct {
		TestField *testFieldType
	}) {
		assert.Equal(t, mrd.str, data.TestField.A)
		didRun = true
	})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	fn(w, r)
	assert.True(t, didRun)
	assert.True(t, c.calledBack)
}

func TestMidwareErrs(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, lerr.Str("Invalid MidwareFunc: func(http.ResponseWriter, *http.Request)"), r)
	}()
	m := New()
	m.Handle(func(w http.ResponseWriter, r *http.Request) {})
}

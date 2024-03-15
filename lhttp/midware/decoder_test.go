package midware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/lhttp/valuedecoder"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name  string
	Age   int
	Admin bool
}

func personFunc(w http.ResponseWriter, r *http.Request, data *struct {
	Form *Person
}) {
	fmt.Fprintf(w, "%+v", data.Form)
}

func TestDecoder(t *testing.T) {
	d := midware.NewDecoder(valuedecoder.Form(), "Form")
	m := midware.New(d)
	h := m.Handle(personFunc)

	form := url.Values{
		"Name":  {"Adam"},
		"Age":   {"39"},
		"Admin": {"false"},
	}
	r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, "&{Name:Adam Age:39 Admin:false}", w.Body.String())
}

// type mockRequestDecoder struct {
// 	str string
// 	err error
// }

// func (mrd *mockRequestDecoder) Decode(i interface{}, r *http.Request) error {
// 	tft := i.(*testFieldType)
// 	tft.A = mrd.str
// 	return mrd.err
// }

// func TestDecoder(t *testing.T) {
// 	mrd := &mockRequestDecoder{
// 		str: "this is a test",
// 	}
// 	d := midware.NewDecoder(mrd, "TestField")

// 	fs := d.Initilize(reflect.TypeOf(testType{}))
// 	tt := &testType{}
// 	fn, err := fs.Insert(nil, nil, reflect.ValueOf(tt))
// 	assert.Nil(t, fn)
// 	assert.NoError(t, err)
// 	assert.Equal(t, mrd.str, tt.TestField.A)
// }

// func TestDecoderNotPtrErr(t *testing.T) {
// 	defer func() {
// 		r := recover()
// 		assert.Equal(t, lerr.Str("Invalid Decoder field: string"), r)
// 	}()
// 	mrd := &mockRequestDecoder{
// 		str: "this is a test",
// 	}
// 	d := NewDecoder(mrd, "A")
// 	d.Initilize(reflect.TypeOf(testFieldType{}))
// }

// func TestDecoderNotStructErr(t *testing.T) {
// 	defer func() {
// 		r := recover()
// 		assert.Equal(t, lerr.Str("Invalid Decoder field: *string"), r)
// 	}()
// 	mrd := &mockRequestDecoder{
// 		str: "this is a test",
// 	}
// 	d := NewDecoder(mrd, "A")
// 	type notStructErr struct {
// 		A *string
// 	}
// 	d.Initilize(reflect.TypeOf(notStructErr{}))
// }

// func TestDecodeErr(t *testing.T) {
// 	mrd := &mockRequestDecoder{
// 		str: "this is a test",
// 		err: lerr.Str("test error"),
// 	}
// 	d := NewDecoder(mrd, "TestField")

// 	fs := d.Initilize(reflect.TypeOf(testType{}))
// 	tt := &testType{}
// 	fn, err := fs.Insert(nil, nil, reflect.ValueOf(tt))
// 	assert.Nil(t, fn)
// 	assert.Equal(t, mrd.err, err)
// 	assert.Nil(t, tt.TestField)
// }

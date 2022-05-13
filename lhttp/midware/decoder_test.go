package midware

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

type mockRequestDecoder struct {
	str string
	err error
}

func (mrd *mockRequestDecoder) Decode(i interface{}, r *http.Request) error {
	tft := i.(*testFieldType)
	tft.A = mrd.str
	return mrd.err
}

func TestDecoder(t *testing.T) {
	mrd := &mockRequestDecoder{
		str: "this is a test",
	}
	d := NewDecoder(mrd, "TestField")

	fs := d.Initilize(reflect.TypeOf(testType{}))
	tt := &testType{}
	fn, err := fs.Insert(nil, nil, reflect.ValueOf(tt))
	assert.Nil(t, fn)
	assert.NoError(t, err)
	assert.Equal(t, mrd.str, tt.TestField.A)
}

func TestDecoderNotPtrErr(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, ErrDecoderField, r)
	}()
	mrd := &mockRequestDecoder{
		str: "this is a test",
	}
	d := NewDecoder(mrd, "A")
	d.Initilize(reflect.TypeOf(testFieldType{}))
}

func TestDecoderNotStructErr(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, ErrDecoderField, r)
	}()
	mrd := &mockRequestDecoder{
		str: "this is a test",
	}
	d := NewDecoder(mrd, "A")
	type notStructErr struct {
		A *string
	}
	d.Initilize(reflect.TypeOf(notStructErr{}))
}

func TestDecodeErr(t *testing.T) {
	mrd := &mockRequestDecoder{
		str: "this is a test",
		err: lerr.Str("test error"),
	}
	d := NewDecoder(mrd, "TestField")

	fs := d.Initilize(reflect.TypeOf(testType{}))
	tt := &testType{}
	fn, err := fs.Insert(nil, nil, reflect.ValueOf(tt))
	assert.Nil(t, fn)
	assert.Equal(t, mrd.err, err)
	assert.Nil(t, tt.TestField)
}

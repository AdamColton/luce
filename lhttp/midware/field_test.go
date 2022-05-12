package midware

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testType struct {
	TestField *testFieldType
}

type mockFieldSetterInitilizer struct {
	isSet bool
}

func (mfsi *mockFieldSetterInitilizer) Initilize(t reflect.Type) FieldSetter {
	return mfsi
}

func (mfsi *mockFieldSetterInitilizer) Set(w http.ResponseWriter, r *http.Request, set reflect.Value) (func(), error) {
	mfsi.isSet = true
	set.Set(reflect.ValueOf(&testFieldType{
		A: "set in Set",
	}))
	return nil, nil
}

func TestFieldBlankErr(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, ErrFieldName, r)
	}()
	f := &FieldInitilizer{}
	f.Initilize(reflect.TypeOf(""))
}

func TestFieldInserterNoField(t *testing.T) {
	f := &FieldInitilizer{
		FieldName: "B",
	}
	di := f.Initilize(reflect.TypeOf(testFieldType{}))
	assert.Nil(t, di)
}

func TestFieldSetter(t *testing.T) {
	mfsi := &mockFieldSetterInitilizer{}
	fi := NewFieldInitilizer(mfsi, "TestField")
	tt := testType{}
	di := fi.Initilize(reflect.TypeOf(tt))
	di.Insert(nil, nil, reflect.ValueOf(&tt))
	assert.True(t, mfsi.isSet)
	assert.Equal(t, "set in Set", tt.TestField.A)
}

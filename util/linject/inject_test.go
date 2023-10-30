package linject_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/linject"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

type StrInit struct {
	Field, Value string
}

func (s StrInit) Initilize(dataType reflect.Type) linject.Injector {
	if linject.CheckField(dataType, s.Field, ltype.String) == nil {
		return nil
	}
	return s
}

func (s StrInit) Inject(data reflect.Value) (func(), error) {
	f := data.Elem().FieldByName(s.Field)
	f.Set(reflect.ValueOf(s.Value))
	return nil, nil
}

func TestFoo(t *testing.T) {
	t.Skip()
	inits := linject.Initilizers{
		StrInit{
			Field: "Name",
			Value: "Adam",
		},
	}

	injs := inits.Initilize(reflector.Type[*Person]())
	assert.NotNil(t, injs)

	p := &Person{}
	cb, err := injs.Inject(reflect.ValueOf(p))
	assert.NoError(t, err)
	assert.Nil(t, cb)
	assert.Equal(t, "Adam", p.Name)
}

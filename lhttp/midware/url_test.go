package midware

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	u := URL("foo", "TestField").(FieldInitilizer).FieldSetterInitilizer.(URLFieldSetter)
	assert.Equal(t, u, u.Initilize(reflect.TypeOf("")))

	restoreVars := Vars
	defer func() {
		Vars = restoreVars
	}()
	Vars = func(r *http.Request) map[string]string {
		return map[string]string{
			"foo": "bar",
		}
	}
	var got struct {
		Foo string
	}
	rg := reflect.ValueOf(&got).Elem().FieldByName("Foo")
	_, err := u.Set(nil, nil, rg)
	assert.NoError(t, err)
	assert.Equal(t, "bar", got.Foo)
}

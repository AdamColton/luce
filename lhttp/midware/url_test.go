package midware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	u := midware.URL("foo", "TestField2")
	m := midware.New(u)

	restoreVars := midware.Vars
	defer func() {
		midware.Vars = restoreVars
	}()
	midware.Vars = func(r *http.Request) map[string]string {
		return map[string]string{
			"foo": "bar",
		}
	}

	r := httptest.NewRequest("Get", "/", nil)
	w := httptest.NewRecorder()
	fn := m.Handle(func(w http.ResponseWriter, r *http.Request, data *struct {
		TestField2 string
	}) {
		w.Write([]byte(data.TestField2))
	})
	fn(w, r)

	assert.Equal(t, w.Body.String(), "bar")
}

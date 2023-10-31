package midware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	m := New()
	m.Initilizer(NewRedirect("Redirect"))
	fn := m.Handle(func(w http.ResponseWriter, r *http.Request, data *struct {
		Redirect string
	}) {
		data.Redirect = "redirectTest"
	})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	fn(w, r)

	assert.Equal(t, "/redirectTest", w.Header().Get("Location"))
}

func TestRedirectErr(t *testing.T) {
	m := New()
	m.Initilizer(NewRedirect("Redirect"))

	defer func() {
		err := recover().(error)
		assert.Equal(t, "expected string, got: int", err.Error())
	}()

	m.Handle(func(w http.ResponseWriter, r *http.Request, data *struct {
		Redirect int
	}) {
		data.Redirect = 123
	})
}

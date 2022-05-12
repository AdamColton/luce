package midware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	m := New(NewRedirect("Redirect"))
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
	m := New(NewRedirect("Redirect"))

	defer func() {
		assert.Equal(t, lerr.Str("Invalid Redirect field: int"), recover())
	}()

	m.Handle(func(w http.ResponseWriter, r *http.Request, data *struct {
		Redirect int
	}) {
		data.Redirect = 123
	})
}

package midware

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

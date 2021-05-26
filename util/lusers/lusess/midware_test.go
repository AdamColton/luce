package lusess_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/util/lusers/lusess"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

// errSaveStore wraps a sessions.Store and always fails Save, to exercise
// Store.Inject's save-failure callback.
type errSaveStore struct{ sessions.Store }

func (s errSaveStore) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	return lerr.Str("Save failed")
}

func userFunc(w http.ResponseWriter, r *http.Request, data *struct {
	Session *lusess.Session
}) {
	u := data.Session.User()
	fmt.Fprint(w, u.Name)
}

func TestMidware(t *testing.T) {
	str := newStore()
	str.FieldName = "Session"

	l := lusess.Login{
		Username: "test-user",
		Password: "test-password",
	}
	_, err := str.Create(l.Username, l.Password)
	assert.NoError(t, err)

	r := loginRequest(l)
	w := httptest.NewRecorder()
	sess, err := str.Login(w, r)
	assert.NoError(t, err)
	sess.Save()

	m := midware.New(str.Midware())
	assert.NotNil(t, m)

	r = httptest.NewRequest("GET", "/", nil)
	r.Header["Cookie"] = w.Header()["Set-Cookie"]
	w = httptest.NewRecorder()
	h := m.Handle(userFunc)
	h(w, r)

	assert.Equal(t, l.Username, w.Body.String())
}

func TestMidwarePanicsOnBlankFieldName(t *testing.T) {
	str := newStore()
	assert.Panics(t, func() {
		str.Midware()
	})
}

func TestMidwareSaveError(t *testing.T) {
	str := newStore()
	str.Store = errSaveStore{str.Store}
	str.FieldName = "Session"

	l := lusess.Login{
		Username: "test-user",
		Password: "test-password",
	}
	_, err := str.Create(l.Username, l.Password)
	assert.NoError(t, err)

	r := loginRequest(l)
	w := httptest.NewRecorder()
	// Login's own Save fails too; that's expected here
	str.Login(w, r)

	m := midware.New(str.Midware())
	r = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	h := m.Handle(func(w http.ResponseWriter, r *http.Request, data *struct {
		Session *lusess.Session
	}) {
		fmt.Fprint(w, "ok")
	})
	// the save failure is only printed, not returned; the handler still runs
	h(w, r)
	assert.Equal(t, "ok", w.Body.String())
}

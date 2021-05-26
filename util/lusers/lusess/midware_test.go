package lusess_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/util/lusers/lusess"
	"github.com/stretchr/testify/assert"
)

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

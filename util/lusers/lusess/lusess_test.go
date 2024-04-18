package lusess_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store/ephemeral/quicknested"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/adamcolton/luce/util/lusers/lusess"
	"github.com/gorilla/schema"
	"github.com/quasoft/memstore"
	"github.com/stretchr/testify/assert"
)

func loginRequest(l lusess.Login) *http.Request {
	form := url.Values{}
	form.Add("Username", l.Username)
	form.Add("Password", l.Password)

	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

func newStore() *lusess.Store {
	keyPairs := [][]byte{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		{17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
	}
	storeFac := quicknested.New(10)
	us := lerr.Must(lusers.NewUserStore(storeFac))
	return &lusess.Store{
		Store:     memstore.NewMemStore(keyPairs...),
		UserStore: us,
		Decoder:   schema.NewDecoder(),
	}
}

func TestLogin(t *testing.T) {
	str := newStore()

	l := lusess.Login{
		Username: "test-user",
		Password: "test-password",
	}
	expected, err := str.Create(l.Username, l.Password)
	assert.NoError(t, err)

	r := loginRequest(l)
	w := httptest.NewRecorder()

	sess, err := str.Login(w, r)
	assert.NoError(t, err)
	assert.NotNil(t, sess)

	got := sess.User()
	assert.Equal(t, expected, got)

	c := w.Result().Cookies()
	assert.Len(t, c, 1)
	assert.Equal(t, "User", c[0].Name)
}

func TestStoreUser(t *testing.T) {
	str := newStore()

	l := lusess.Login{
		Username: "test-user",
		Password: "test-password",
	}
	expected, err := str.Create(l.Username, l.Password)
	assert.NoError(t, err)

	// no session yet: no user, no error
	fresh := httptest.NewRequest(http.MethodGet, "/", nil)
	u, err := str.User(fresh)
	assert.NoError(t, err)
	assert.Nil(t, u)

	r := loginRequest(l)
	w := httptest.NewRecorder()
	_, err = str.Login(w, r)
	assert.NoError(t, err)

	// replay the session cookie Login set
	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range w.Result().Cookies() {
		r2.AddCookie(c)
	}
	u, err = str.User(r2)
	assert.NoError(t, err)
	assert.Equal(t, expected, u)
}

func TestStoreSessionDecodeError(t *testing.T) {
	str := newStore()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: lusess.StoreName, Value: "not a valid cookie"})

	_, err := str.Session(nil, r)
	assert.Error(t, err)

	_, err = str.User(r)
	assert.Error(t, err)
}

func TestLoginFailure(t *testing.T) {
	str := newStore()

	l := lusess.Login{
		Username: "test-user",
		Password: "test-password",
	}
	_, err := str.Create(l.Username, l.Password)
	assert.NoError(t, err)

	wrong := lusess.Login{
		Username: "test-user",
		Password: "not-the-password",
	}
	r := loginRequest(wrong)
	w := httptest.NewRecorder()

	sess, err := str.Login(w, r)
	assert.ErrorIs(t, err, lusess.ErrLoginFailed)
	assert.Nil(t, sess.User())
	// a failed login must not set a session cookie
	assert.Empty(t, w.Result().Cookies())
}

func TestLoginDecodeError(t *testing.T) {
	str := newStore()

	form := url.Values{}
	form.Add("Username", "test-user")
	form.Add("Password", "test-password")
	form.Add("NotAField", "x")
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := str.Login(httptest.NewRecorder(), r)
	assert.Error(t, err)
}

func TestLoginSessionError(t *testing.T) {
	str := newStore()

	l := lusess.Login{
		Username: "test-user",
		Password: "test-password",
	}
	_, err := str.Create(l.Username, l.Password)
	assert.NoError(t, err)

	r := loginRequest(l)
	r.AddCookie(&http.Cookie{Name: lusess.StoreName, Value: "not a valid cookie"})

	_, err = str.Login(httptest.NewRecorder(), r)
	assert.Error(t, err)
}

func TestSessionSetUserAndLogout(t *testing.T) {
	str := newStore()
	u, err := str.Create("test-user", "test-password")
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	sess, err := str.Session(w, r)
	assert.NoError(t, err)
	assert.Nil(t, sess.User())

	sess.SetUser(u)
	assert.Equal(t, u, sess.User())

	sess.Logout()
	assert.Nil(t, sess.User())
}

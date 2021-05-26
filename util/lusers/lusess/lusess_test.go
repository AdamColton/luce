package lusess_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store/ephemeral"
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
	storeFac := ephemeral.Factory(bytebtree.New, 10)
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

	assert.NoError(t, sess.Save())
	c := w.Result().Cookies()
	assert.Len(t, c, 1)
	assert.Equal(t, "User", c[0].Name)
}

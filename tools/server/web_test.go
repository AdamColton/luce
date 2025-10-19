package server_test

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestServerWebFlow(t *testing.T) {
	cfg := testConfig(":54420")
	srv, err := (&cfg).New()
	assert.NoError(t, err)

	_, err = srv.Users.Create("bob", "secret")
	assert.NoError(t, err)

	closed := make(chan bool)
	go func() {
		srv.Run()
		closed <- true
	}()
	defer func() {
		assert.NoError(t, srv.Close())
		timeout.After(300, closed)
	}()

	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)
	client := &http.Client{Jar: jar, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	base := "http://localhost:54420"

	var resp *http.Response
	err = timeout.After(500, func() {
		for {
			resp, err = client.Get(base + "/")
			if err == nil {
				return
			}
		}
	})
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "home", string(body))

	resp, err = client.Get(base + "/user/signin")
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "signin")

	form := url.Values{"Username": {"bob"}, "Password": {"secret"}}
	resp, err = client.PostForm(base+"/user/signin", form)
	assert.NoError(t, err)
	io.ReadAll(resp.Body)

	// signed in now, so GET /user/signin should redirect away
	resp, err = client.Get(base + "/user/signin")
	assert.NoError(t, err)
	io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusFound, resp.StatusCode)

	resp, err = client.Get(base + "/")
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "home, signed in", string(body))

	resp, err = client.Get(base + "/services")
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Services")

	resp, err = client.Get(base + "/admin/users")
	assert.NoError(t, err)
	io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	loc, err := resp.Location()
	assert.NoError(t, err)
	assert.Equal(t, "/", loc.Path)

	resp, err = client.Get(base + "/user/signout")
	assert.NoError(t, err)
	io.ReadAll(resp.Body)

	resp, err = client.Get(base + "/")
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "home", string(body))

	badForm := url.Values{"Username": {"bob"}, "Password": {"wrong"}}
	resp, err = client.PostForm(base+"/user/signin", badForm)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	loc, err = resp.Location()
	assert.NoError(t, err)
	assert.Contains(t, loc.String(), "login+failed")
}

func TestServerWebAdminLoggedIn(t *testing.T) {
	cfg := testConfig(":54421")
	srv, err := (&cfg).New()
	assert.NoError(t, err)
	u, err := srv.Users.Create("admin", "secret")
	assert.NoError(t, err)
	g, err := srv.Users.Group("admin")
	assert.NoError(t, err)
	assert.NoError(t, g.AddUser(u))

	closed := make(chan bool)
	go func() {
		srv.Run()
		closed <- true
	}()
	defer func() {
		assert.NoError(t, srv.Close())
		timeout.After(300, closed)
	}()

	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)
	client := &http.Client{Jar: jar}
	base := "http://localhost:54421"

	err = timeout.After(500, func() {
		for {
			_, e := client.Get(base + "/")
			if e == nil {
				return
			}
		}
	})
	assert.NoError(t, err)

	form := url.Values{"Username": {"admin"}, "Password": {"secret"}}
	resp, err := client.PostForm(base+"/user/signin", form)
	assert.NoError(t, err)
	io.ReadAll(resp.Body)

	resp, err = client.Get(base + "/admin/users")
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "admin")

	resp, err = client.Get(base + "/admin/listBashCmds")
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "bash commands")
}

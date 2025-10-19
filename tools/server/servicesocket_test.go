package server_test

import (
	"io"
	"net"
	"net/http"
	"testing"

	servicepkg "github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// TestServiceSocketProxy exercises the full service-registration and
// request-proxying flow: a service connects over the unix service socket,
// registers a route, and an HTTP request against the main server gets
// proxied to it and back.
func TestServiceSocketProxy(t *testing.T) {
	sock := t.TempDir() + "/service.sock"
	cfg := testConfig(":54430")
	cfg.ServiceSocket = sock
	srv, err := (&cfg).New()
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

	err = timeout.After(500, func() {
		for {
			conn, e := net.Dial("unix", sock)
			if e == nil {
				conn.Close()
				return
			}
		}
	})
	assert.NoError(t, err)

	client, err := servicepkg.NewService("widget", "", "/widget", sock)
	assert.NoError(t, err)
	client.Service.AddLink("hello", "", "hello")

	route := servicepkg.NewRoute("hello").Get()
	client.Add(func(r *servicepkg.Request) *servicepkg.Response {
		return r.ResponseString("hello from widget")
	}, route)

	go client.Run()

	var resp *http.Response
	err = timeout.After(1000, func() {
		for {
			resp, err = http.Get("http://localhost:54430/widget/hello")
			// the route isn't registered until the service's registration
			// message has been processed over the socket, so a 404 means
			// "not ready yet", not "no such route" -- keep retrying.
			if err == nil && resp.StatusCode != http.StatusNotFound {
				return
			}
		}
	})
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "hello from widget", string(body))

	resp, err = http.Get("http://localhost:54430/services")
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "widget")
	assert.Contains(t, string(body), "hello")
}

// TestServiceSocketReconnect covers serviceRoute.setActive: a route goes
// inactive when its service disconnects, and active again when the same
// route is registered a second time.
func TestServiceSocketReconnect(t *testing.T) {
	sock := t.TempDir() + "/service2.sock"
	cfg := testConfig(":54431")
	cfg.ServiceSocket = sock
	srv, err := (&cfg).New()
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

	err = timeout.After(500, func() {
		for {
			conn, e := net.Dial("unix", sock)
			if e == nil {
				conn.Close()
				return
			}
		}
	})
	assert.NoError(t, err)

	connect := func() *servicepkg.Client {
		client, err := servicepkg.NewService("widget2", "", "/widget2", sock)
		assert.NoError(t, err)
		route := servicepkg.NewRoute("hello").Get()
		client.Add(func(r *servicepkg.Request) *servicepkg.Response {
			return r.ResponseString("hello again")
		}, route)
		go client.Run()
		return client
	}

	c1 := connect()
	var resp *http.Response
	err = timeout.After(1000, func() {
		for {
			resp, err = http.Get("http://localhost:54431/widget2/hello")
			if err == nil && resp.StatusCode != http.StatusNotFound {
				return
			}
		}
	})
	assert.NoError(t, err)
	io.ReadAll(resp.Body)

	// disconnecting makes the route inactive
	c1.NetConn.Close()
	err = timeout.After(1000, func() {
		for {
			resp, err = http.Get("http://localhost:54431/widget2/hello")
			if err != nil || resp.StatusCode == http.StatusNotFound {
				return
			}
		}
	})
	assert.NoError(t, err)

	// reconnecting with the same route makes it active again
	connect()
	err = timeout.After(1000, func() {
		for {
			resp, err = http.Get("http://localhost:54431/widget2/hello")
			if err == nil && resp.StatusCode == http.StatusOK {
				return
			}
		}
	})
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "hello again", string(body))
}

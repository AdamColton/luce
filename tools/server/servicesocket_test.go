package server_test

import (
	"io"
	"net"
	"net/http"
	"testing"

	servicepkg "github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/gorilla/websocket"
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

// TestServiceSocketReconnectMoved covers a service that comes back with its
// route somewhere new, as when a service's Base is changed while the server
// keeps running. The route has to answer at the new path, not the old one.
func TestServiceSocketReconnectMoved(t *testing.T) {
	sock := t.TempDir() + "/service.sock"
	cfg := testConfig(":54433")
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

	connect := func(base string) *servicepkg.Client {
		client, err := servicepkg.NewService("widget", "", base, sock)
		assert.NoError(t, err)
		route := servicepkg.NewRoute("hello").Get()
		client.Add(func(r *servicepkg.Request) *servicepkg.Response {
			return r.ResponseString("hello from " + base)
		}, route)
		go client.Run()
		return client
	}

	// awaitStatus polls url until it answers with status, since registering
	// and disconnecting both take effect asynchronously over the socket.
	awaitStatus := func(url string, status int) *http.Response {
		var resp *http.Response
		err := timeout.After(1000, func() {
			for {
				r, err := http.Get(url)
				if err == nil && r.StatusCode == status {
					resp = r
					return
				}
			}
		})
		assert.NoError(t, err)
		return resp
	}

	const (
		oldHello = "http://localhost:54433/widget/hello"
		newHello = "http://localhost:54433/gadget/hello"
	)

	c := connect("/widget")
	awaitStatus(oldHello, http.StatusOK)
	assert.NoError(t, c.NetConn.Close())
	awaitStatus(oldHello, http.StatusNotFound)

	// Back under a new Base: the route moves with it.
	connect("/gadget")
	resp := awaitStatus(newHello, http.StatusOK)
	if resp != nil {
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "hello from /gadget", string(body))
	}

	resp, err = http.Get(oldHello)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestServiceSocketSharedRoutePath covers two services that each declare a
// route at the same path under their own Base, as journal.service and
// stt.service both declare "home". When one disconnects, only its own route
// goes inactive; the other service keeps serving.
func TestServiceSocketSharedRoutePath(t *testing.T) {
	sock := t.TempDir() + "/service.sock"
	cfg := testConfig(":54432")
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

	connect := func(name, base string) *servicepkg.Client {
		client, err := servicepkg.NewService(name, "", base, sock)
		assert.NoError(t, err)
		route := servicepkg.NewRoute("home").Get()
		client.Add(func(r *servicepkg.Request) *servicepkg.Response {
			return r.ResponseString(name + " home")
		}, route)
		go client.Run()
		return client
	}

	// awaitStatus polls url until it answers with status, since registering
	// and disconnecting both take effect asynchronously over the socket.
	awaitStatus := func(url string, status int) *http.Response {
		var resp *http.Response
		err := timeout.After(1000, func() {
			for {
				r, err := http.Get(url)
				if err == nil && r.StatusCode == status {
					resp = r
					return
				}
			}
		})
		assert.NoError(t, err)
		return resp
	}

	const (
		journalHome = "http://localhost:54432/journal.service/home"
		sttHome     = "http://localhost:54432/stt.service/home"
	)

	// Journal registers before stt, so they register in a known order.
	journal := connect("journal", "/journal.service")
	awaitStatus(journalHome, http.StatusOK)
	connect("stt", "/stt.service")
	awaitStatus(sttHome, http.StatusOK)

	// Journal restarting takes down journal's home and nothing else.
	assert.NoError(t, journal.NetConn.Close())
	awaitStatus(journalHome, http.StatusNotFound)

	resp, err := http.Get(sttHome)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "stt home", string(body))
}

// TestServiceDisconnectClosesSockets covers what happens to a client's
// websocket when the service behind it disconnects, as a restart does. The
// server owns the socket, so nothing at the client's end can tell the service
// has gone: unless the server closes it, the client holds a connection that
// looks alive and never reconnects.
func TestServiceDisconnectClosesSockets(t *testing.T) {
	sock := t.TempDir() + "/service.sock"
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

	svc, err := servicepkg.NewService("widget", "", "/widget", sock)
	assert.NoError(t, err)

	opened := make(chan uint32, 1)
	servicepkg.NewRoute("ws").
		WithSocket().
		Handle(svc, func(r *servicepkg.Request) *servicepkg.Response {
			opened <- r.WebSocket
			return r.Response(nil)
		})
	go svc.Run()

	// The route isn't live until the registration has been processed, so a
	// failed dial means "not ready yet".
	var ws *websocket.Conn
	err = timeout.After(1000, func() {
		for {
			ws, _, err = websocket.DefaultDialer.Dial("ws://localhost:54431/widget/ws", nil)
			if err == nil {
				return
			}
		}
	})
	assert.NoError(t, err)
	defer ws.Close()
	assert.NoError(t, timeout.After(1000, func() { <-opened }))

	// The service goes away.
	assert.NoError(t, svc.NetConn.Close())

	// The client's socket has to close, because that is the only signal it
	// can act on. Without it this read blocks for good.
	assert.NoError(t, timeout.After(1000, func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				return
			}
		}
	}))
}

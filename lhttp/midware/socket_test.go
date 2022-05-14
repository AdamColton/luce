package midware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSocket(t *testing.T) {
	closed := make(chan bool)
	fn := NewWebSocket().HandleSocketChans(func(to chan<- []byte, from <-chan []byte, r *http.Request) {
		for m := range from {
			m[1] += m[0]
			m[0]++
			to <- m
		}
		closed <- true
	})

	s := httptest.NewServer(http.HandlerFunc(fn))
	defer s.Close()

	u := strings.Replace(s.URL, "http", "ws", 1)
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)

	m := []byte{1, 2}
	conn.WriteMessage(websocket.BinaryMessage, m)
	_, m, err = conn.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, []byte{2, 3}, m)
	conn.WriteMessage(websocket.BinaryMessage, m)
	_, m, err = conn.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, []byte{3, 5}, m)

	conn.Close()

	timeout.After(100, closed)
}

type magicFinalizer struct {
	ch chan bool
}

func (mf *magicFinalizer) Initilize(t reflect.Type) DataInserter {
	return mf
}

func (mf *magicFinalizer) Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) (func(), error) {
	mf.ch = make(chan bool)
	return func() {
		mf.ch <- true
	}, nil
}

func TestWebSocketChannelInitilizer(t *testing.T) {
	ws := NewWebSocket()
	mf := &magicFinalizer{}
	m := New(ws.Initilizer("To", "From", ""), mf)

	shouldclose := make(chan bool)
	s := httptest.NewServer(m.Handle(func(w http.ResponseWriter, r *http.Request, data struct {
		To   chan<- []byte
		From <-chan []byte
	}) {
		err := timeout.After(100, func() {
			assert.Equal(t, "client to server", string(<-data.From))
		})
		assert.NoError(t, err)
		data.To <- []byte("server to client")
		<-shouldclose
		close(data.To)
	}))
	defer s.Close()

	u := strings.Replace(s.URL, "http", "ws", 1)
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	err = timeout.After(100, func() error {
		return conn.WriteMessage(websocket.BinaryMessage, []byte("client to server"))
	})
	assert.NoError(t, err)

	err = timeout.After(100, func() error {
		_, got, err := conn.ReadMessage()
		assert.Equal(t, []byte("server to client"), got)
		return err
	})
	assert.NoError(t, err)
	shouldclose <- true
	<-mf.ch
}

func TestWebSocketInitilizer(t *testing.T) {
	ws := NewWebSocket()
	m := New(ws.Initilizer("", "", "Socket"))

	shouldclose := make(chan bool)
	s := httptest.NewServer(m.Handle(func(w http.ResponseWriter, r *http.Request, data struct {
		Socket *websocket.Conn
	}) {
		err := timeout.After(100, func() error {
			_, msg, err := data.Socket.ReadMessage()
			assert.Equal(t, "client to server", string(msg))
			return err
		})
		assert.NoError(t, err)
		err = data.Socket.WriteMessage(websocket.BinaryMessage, []byte("server to client"))
		assert.NoError(t, err)
		<-shouldclose
	}))
	defer s.Close()

	u := strings.Replace(s.URL, "http", "ws", 1)
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	err = timeout.After(100, func() error {
		return conn.WriteMessage(websocket.BinaryMessage, []byte("client to server"))
	})
	assert.NoError(t, err)

	err = timeout.After(100, func() error {
		_, got, err := conn.ReadMessage()
		assert.Equal(t, []byte("server to client"), got)
		return err
	})
	assert.NoError(t, err)
	shouldclose <- true
}

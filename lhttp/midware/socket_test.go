package midware

import (
	"net/http"
	"net/http/httptest"
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

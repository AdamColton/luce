package midware

import (
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp"
	"github.com/gorilla/websocket"
)

type WebSocket struct {
	Upgrader       websocket.Upgrader
	ToBuf, FromBuf int
	lhttp.ErrHandler
}

func NewWebSocket() WebSocket {
	return WebSocket{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (ws WebSocket) Handler(handler lhttp.SocketHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		socket, err := ws.Upgrader.Upgrade(w, r, nil)
		if !ws.Check(w, r, lerr.Wrap(err, "while_upgrading_socket")) {
			handler(socket, r)
		}
	}
}

// HandleSocketChans abstracts the websocket as a pair of channels. The handler
// must close the to channel when it is done.
func (ws WebSocket) HandleSocketChans(handler lhttp.ChanHandler) http.HandlerFunc {
	return ws.Handler(func(socket *websocket.Conn, r *http.Request) {
		to := make(chan []byte, ws.ToBuf)
		from := make(chan []byte, ws.FromBuf)

		socket.SetCloseHandler(func(code int, text string) error {
			close(to)
			return nil
		})

		sw := lhttp.NewSocket(socket)

		go sw.RunReader(from)
		go handler(to, from, r)
		sw.RunSender(to)
		socket.Close()
	})
}

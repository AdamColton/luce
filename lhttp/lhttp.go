package lhttp

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// RequestDecoder encapsulates decoding a request to a value
type RequestDecoder interface {
	Decode(interface{}, *http.Request) error
}

// SocketHandler is similar to http.HandlerFunc, but for handling websockets
type SocketHandler func(*websocket.Conn, *http.Request)

// ChanHandler represents duplex communication with two channels.
type ChanHandler func(to chan<- []byte, from <-chan []byte, r *http.Request)

// ErrHandler can be used to
type ErrHandler func(w http.ResponseWriter, r *http.Request, err error)

// Check if err is not nill and call underlying ErrHandler if it is not nil
func (h ErrHandler) Check(w http.ResponseWriter, r *http.Request, err error) bool {
	isErr := err != nil
	if isErr && h != nil {
		h(w, r, err)
	}
	return isErr
}

type MessageReaderWriter interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (int, []byte, error)
}

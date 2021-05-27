package lhttp

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// RequestDecoder decodes the data of a request into a value.
type RequestDecoder interface {
	// Decode fills dst, which must be a pointer, from r.
	Decode(interface{}, *http.Request) error
}

// SocketHandler is similar to http.HandlerFunc, but it is given a websocket
// connection as well as the request.
type SocketHandler func(*websocket.Conn, *http.Request)

// ChanHandler handles a duplex connection as two channels. It sends messages to
// the client on to and receives the client's messages from from.
type ChanHandler func(to chan<- []byte, from <-chan []byte, r *http.Request)

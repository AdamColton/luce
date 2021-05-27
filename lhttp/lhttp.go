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

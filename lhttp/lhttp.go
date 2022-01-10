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

// StatusCoder allows errors to indicate the status code they should return.
type StatusCoder interface {
	StatusCode() int
}

// StatusCodeWrapper combines an error and a status code.
type StatusCodeWrapper struct {
	Err    error
	Status int
}

// Error fullfils error.
func (scw StatusCodeWrapper) Error() string {
	return scw.Err.Error()
}

// StatusCode fulfills StatusCoder.
func (scw StatusCodeWrapper) StatusCode() int {
	return scw.Status
}

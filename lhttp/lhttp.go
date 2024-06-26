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

// ErrHandler is intended to repsond to a request when an error has occured.
type ErrHandler func(w http.ResponseWriter, r *http.Request, err error)

// Check only invokes the ErrHandler if err is not nil. Returns a bool
// indicating if ErrHandler was invoked.
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

// StatusErr provides an optional interface for errors allowing errors to
// include HTTP status.
type StatusErr interface {
	Status() int
}

// ErrStatus will return 0 if err is nil. If err is not nil and fulfills
// StatusErr, that Status will be returned, otherwise StatusInternalServerError
// is returned.
func ErrStatus(err error) int {
	if err == nil {
		return 0
	}
	if s, ok := err.(StatusErr); ok {
		return s.Status()
	}
	return http.StatusInternalServerError
}

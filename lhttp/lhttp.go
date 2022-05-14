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

// ErrHandler responds to a request when an error has occurred. A nil ErrHandler
// is valid: Check then only reports whether there was an error.
type ErrHandler func(w http.ResponseWriter, r *http.Request, err error)

// Check invokes the ErrHandler if err is not nil and the ErrHandler is not nil.
// It returns true if err is not nil, whether or not the ErrHandler was invoked.
func (h ErrHandler) Check(w http.ResponseWriter, r *http.Request, err error) bool {
	isErr := err != nil
	if isErr && h != nil {
		h(w, r, err)
	}
	return isErr
}

// MessageReaderWriter is the part of a websocket connection (such as
// *websocket.Conn) that Socket needs: reading and writing whole messages. The
// messageType is a websocket message type such as websocket.TextMessage.
type MessageReaderWriter interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (int, []byte, error)
}

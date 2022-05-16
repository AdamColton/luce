package service

import (
	"encoding/gob"
	"net/url"

	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/lusers"
)

var tm = type32.NewTypeMap()

func init() {
	gob.Register((Routes)(nil))
	gob.Register(Request{})
	gob.Register(Response{})
	gob.Register(SocketOpened{})
	gob.Register(SocketClose{})
	gob.Register(SocketMessage{})

	tm.RegisterType32s(
		(Routes)(nil),
		Request{},
		Response{},
		SocketOpened{},
		SocketClose{},
		SocketMessage{},
	)
}

// Request represents a user request that the luce Server is relaying to the
// service.
type Request struct {
	ID          uint32
	RouteConfig string
	Path        string
	Method      string
	PathVars    map[string]string
	Query       map[string]string
	Form        url.Values
	Body        []byte
	User        *lusers.User
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (Request) TypeID32() uint32 {
	return 161709784
}

// Response to the Request.
func (r Request) Response(body []byte) Response {
	return Response{
		ID:   r.ID,
		Body: body,
	}
}

// ResponseString to the Request.
func (r Request) ResponseString(body string) Response {
	return r.Response([]byte(body))
}

// Response to a request. The ID is the same as the ID is taken from the
// request.
type Response struct {
	ID   uint32
	Body []byte
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (Response) TypeID32() uint32 {
	return 370114636
}

type SocketOpened struct {
	ID uint32
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketOpened) TypeID32() uint32 {
	return 1046109042
}

type SocketClose struct {
	ID uint32
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketClose) TypeID32() uint32 {
	return 3196974518
}

type SocketMessage struct {
	ID   uint32
	Body []byte
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketMessage) TypeID32() uint32 {
	return 3196974518
}

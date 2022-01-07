package service

import (
	"encoding/gob"
	"net/http"
	"net/url"

	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/lusers"
)

var tm = type32.NewTypeMap()

func Register(zeroValues ...type32.TypeIDer32) {
	for _, z := range zeroValues {
		gob.Register(z)
	}
	tm.RegisterType32s(zeroValues...)
}

func init() {
	Register(
		(Routes)(nil),
		Request{},
		Response{},
		SocketOpened{},
		SocketClose{},
		SocketMessage{},
	)
}

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

func (Request) TypeID32() uint32 {
	return 161709784
}

func (r Request) Response(body []byte) Response {
	return Response{
		ID:     r.ID,
		Body:   body,
		Status: http.StatusOK,
	}
}

func (r Request) ResponseString(body string) Response {
	return r.Response([]byte(body))
}

type Response struct {
	ID     uint32
	Body   []byte
	Status int
}

func (Response) TypeID32() uint32 {
	return 370114636
}

func (r Response) SetStatus(status int) Response {
	r.Status = status
	return r
}

type SocketOpened struct {
	ID uint32
}

func (SocketOpened) TypeID32() uint32 {
	return 1046109042
}

type SocketClose struct {
	ID uint32
}

func (SocketClose) TypeID32() uint32 {
	return 3196974518
}

type SocketMessage struct {
	ID   uint32
	Body []byte
}

func (SocketMessage) TypeID32() uint32 {
	return 3196974518
}

package service

import (
	"encoding/gob"
	"net/http"
	"net/url"

	"github.com/adamcolton/luce/lhttp"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/lusers"
)

var tm = type32.NewTypeMap()

// Register types with both gob and the typemap.
func Register(zeroValues ...type32.TypeIDer32) {
	for _, z := range zeroValues {
		gob.Register(z)
	}
	tm.RegisterType32s(zeroValues...)
}

func init() {
	Register(
		(Routes)(nil),
		(*Request)(nil),
		(*Response)(nil),
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
func (*Request) TypeID32() uint32 {
	return 161709784
}

// Response to the Request.
func (r *Request) Response(body []byte) *Response {
	return &Response{
		ID:     r.ID,
		Body:   body,
		Status: http.StatusOK,
	}
}

// ResponseString to the Request.
func (r *Request) ResponseString(body string) *Response {
	return r.Response([]byte(body))
}

// ResponseError to the Request. If status==0 and err implements
// lhttp.StatusCoder, that status code will be used.
func (r *Request) ResponseError(err error, status int) *Response {
	var resp *Response
	if err == nil {
		resp = r.Response(nil)
	} else {
		resp = r.ResponseString(err.Error())
	}
	if status == 0 {
		if sc, ok := err.(lhttp.StatusCoder); ok {
			resp.SetStatus(sc.StatusCode())
		}
	}
	return resp
}

// Response to a request. The ID is the same as the ID is taken from the
// request.
type Response struct {
	ID     uint32
	Body   []byte
	Status int
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (*Response) TypeID32() uint32 {
	return 370114636
}

func (r *Response) SetStatus(status int) *Response {
	r.Status = status
	return r
}

func (r *Response) Write(p []byte) (n int, err error) {
	if r.Body == nil {
		r.Body = p
	} else {
		r.Body = append(r.Body, p...)
	}
	return len(p), nil
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

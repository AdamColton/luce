package service

import (
	"net/http"
	"net/url"

	"github.com/adamcolton/luce/util/lusers"
)

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

const RequestTypeID32 uint32 = 161709784

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (*Request) TypeID32() uint32 {
	return RequestTypeID32
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

// ResponseErr sets the response body to the error and sets the status.
func (r *Request) ResponseErr(err error, status int) *Response {
	resp := r.ResponseString(err.Error())
	resp.Status = status
	return resp
}

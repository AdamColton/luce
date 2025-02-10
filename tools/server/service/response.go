package service

import (
	"net/http"
)

// Response to a request. The ID is the same as the ID is taken from the
// request.
type Response struct {
	ID     uint32
	Body   []byte
	Status int
	Header http.Header
}

const ResponseTypeID32 uint32 = 370114636

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (*Response) TypeID32() uint32 {
	return ResponseTypeID32
}

// Set Response HTTP Status Code
func (r *Response) SetStatus(status int) *Response {
	r.Status = status
	return r
}

func (r *Response) SetHeader(key, val string) *Response {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Set(key, val)
	return r
}

const ContentType = "Content-Type"

func (r *Response) ContentType(val string) *Response {
	return r.SetHeader(ContentType, val)
}

// Write p to the body. Fulfills io.Writer.
func (r *Response) Write(p []byte) (n int, err error) {
	if r.Body == nil {
		r.Body = p
	} else {
		r.Body = append(r.Body, p...)
	}
	return len(p), nil
}

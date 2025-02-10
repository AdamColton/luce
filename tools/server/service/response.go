package service

import "net/http"

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

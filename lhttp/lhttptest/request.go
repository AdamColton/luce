// Package lhttptest has helpers for building the *http.Request values that HTTP
// handlers and decoders are tested with.
package lhttptest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

// Request describes a request to build with GET or POST.
type Request struct {
	// Target is the path of the request, used by GET.
	Target string
	// Values are sent as the query string by GET and as the form body by POST.
	Values url.Values
}

// NewRequest creates a Request for target. src is nil (no values), a
// url.Values (used as is) or a struct, which is encoded with gorilla/schema. If
// a src cannot be encoded, Values is left empty and the error is dropped; call
// EncodeValues to see it.
func NewRequest(target string, src any) *Request {
	r := &Request{
		Target: target,
	}
	if src != nil {
		if uv, ok := src.(url.Values); ok {
			r.Values = uv
		} else {
			r.EncodeValues(src)
		}
	}
	return r
}

// POST returns a POST request to "/" (Target is not used) whose body is Values
// as a form. Without Values it has no body and no content type.
func (req *Request) POST() *http.Request {
	var rdr io.Reader
	hasForm := req.Values != nil
	if hasForm {
		rdr = strings.NewReader(req.Values.Encode())
	}
	r := httptest.NewRequest("POST", "/", rdr)
	if hasForm {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return r
}

// GET returns a GET request to Target, with Values as the query string.
func (req *Request) GET() *http.Request {
	t := req.Target
	if req.Values != nil {
		t = t + "?" + req.Values.Encode()
	}
	return httptest.NewRequest("GET", t, nil)
}

// Encoder is the gorilla/schema encoder EncodeValues uses. Get creates it the
// first time it is needed.
var Encoder *schema.Encoder

// Get returns Encoder, creating it if it is nil.
func Get() *schema.Encoder {
	if Encoder == nil {
		Encoder = schema.NewEncoder()
	}
	return Encoder
}

// EncodeValues replaces Values with the encoding of src, which must be a
// struct, and returns the error from the encoder.
func (req *Request) EncodeValues(src any) error {
	enc := Get()
	req.Values = make(url.Values)
	return enc.Encode(src, req.Values)
}

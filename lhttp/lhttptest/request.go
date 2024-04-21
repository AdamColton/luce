package lhttptest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

type Request struct {
	Target string
	Values url.Values
}

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

func (req *Request) GET() *http.Request {
	t := req.Target
	if req.Values != nil {
		t = t + "?" + req.Values.Encode()
	}
	return httptest.NewRequest("GET", t, nil)
}

var Encoder *schema.Encoder

func Get() *schema.Encoder {
	if Encoder == nil {
		Encoder = schema.NewEncoder()
	}
	return Encoder
}

func (req *Request) EncodeValues(src any) error {
	enc := Get()
	req.Values = make(url.Values)
	return enc.Encode(src, req.Values)
}

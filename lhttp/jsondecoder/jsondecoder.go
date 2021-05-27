package jsondecoder

import (
	"encoding/json"
	"net/http"
)

// JsonDecoder fulfills lhttp.RequestDecoder
type JsonDecoder struct{}

// New JsonDecoder
func New() JsonDecoder {
	return JsonDecoder{}
}

// Decode the request body as json into dst.
func (jd JsonDecoder) Decode(dst interface{}, r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

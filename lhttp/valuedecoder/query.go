package valuedecoder

import (
	"net/http"
	"net/url"
)

// Query returns a Decoder for the query string of a request.
func Query() Decoder {
	return Decoder{
		URLValuesDecoder: Get(),
		Getter:           QueryGetter,
	}
}

// QueryGetter returns the query values of the request. It never fails.
func QueryGetter(r *http.Request) (url.Values, error) {
	return r.URL.Query(), nil
}

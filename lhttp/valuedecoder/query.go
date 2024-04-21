package valuedecoder

import (
	"net/http"
	"net/url"
)

func Query() Decoder {
	return Decoder{
		URLValuesDecoder: Get(),
		Getter:           QueryGetter,
	}
}

func QueryGetter(r *http.Request) (url.Values, error) {
	return r.URL.Query(), nil
}

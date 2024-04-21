package valuedecoder

import (
	"net/http"
	"net/url"
)

func Form() Decoder {
	return Decoder{
		URLValuesDecoder: Get(),
		Getter:           FormGetter,
	}
}

func FormGetter(r *http.Request) (url.Values, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	return r.PostForm, nil
}

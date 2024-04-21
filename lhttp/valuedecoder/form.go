package valuedecoder

import (
	"net/http"
	"net/url"
)

// Form returns a Decoder for the form body of a POST request.
func Form() Decoder {
	return Decoder{
		URLValuesDecoder: Get(),
		Getter:           FormGetter,
	}
}

// FormGetter parses the request and returns its POST form values.
func FormGetter(r *http.Request) (url.Values, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	return r.PostForm, nil
}

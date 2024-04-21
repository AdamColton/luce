package valuedecoder

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

type URLValuesDecoder interface {
	Decode(interface{}, map[string][]string) error
}

// Shared is exported so that it can be reused.
var Shared URLValuesDecoder

func Get() URLValuesDecoder {
	if Shared == nil {
		Shared = schema.NewDecoder()
	}
	return Shared
}

type Decoder struct {
	Getter func(*http.Request) (url.Values, error)
	URLValuesDecoder
}

func (d Decoder) Decode(dst interface{}, r *http.Request) error {
	data, err := d.Getter(r)
	if err != nil {
		return err
	}
	return d.URLValuesDecoder.Decode(dst, data)
}

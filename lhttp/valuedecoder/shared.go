package valuedecoder

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

// URLValuesDecoder fills a value from url values. *schema.Decoder fulfills it.
type URLValuesDecoder interface {
	Decode(interface{}, map[string][]string) error
}

// Shared is exported so that it can be reused.
var Shared URLValuesDecoder

// Get returns Shared, creating a gorilla/schema decoder the first time.
func Get() URLValuesDecoder {
	if Shared == nil {
		Shared = schema.NewDecoder()
	}
	return Shared
}

// Decoder gets the values out of a request and decodes them. Getter must be
// set: Decode panics on a nil Getter.
type Decoder struct {
	// Getter reads the values to decode out of the request.
	Getter func(*http.Request) (url.Values, error)
	URLValuesDecoder
}

// Decode fills dst from the values Getter returns. It fulfills
// lhttp.RequestDecoder.
func (d Decoder) Decode(dst interface{}, r *http.Request) error {
	data, err := d.Getter(r)
	if err != nil {
		return err
	}
	return d.URLValuesDecoder.Decode(dst, data)
}

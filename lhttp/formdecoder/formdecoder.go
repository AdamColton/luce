package formdecoder

import (
	"net/http"

	"github.com/gorilla/schema"
)

// SchemaDecoder is exported so that it can be reused.
var SchemaDecoder *schema.Decoder

// New FormDecoder, creates a form decoder using the default SchemaDecoder. If
// SchemaDecoder is nil, it is initilized.
func New() FormDecoder {
	if SchemaDecoder == nil {
		SchemaDecoder = schema.NewDecoder()
	}
	return FormDecoder{
		SchemaDecoder: SchemaDecoder,
	}
}

// FormDecoder fulfills lhttp.RequestDecoder
type FormDecoder struct {
	SchemaDecoder interface {
		Decode(interface{}, map[string][]string) error
	}
}

// Decode a form submission from request r into dst.
func (fd FormDecoder) Decode(dst interface{}, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	return fd.SchemaDecoder.Decode(dst, r.PostForm)
}

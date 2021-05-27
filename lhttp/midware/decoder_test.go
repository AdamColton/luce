package midware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/adamcolton/luce/lhttp/formdecoder"
	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name  string
	Age   int
	Admin bool
}

func personFunc(w http.ResponseWriter, r *http.Request, data struct {
	Form *Person
}) {
	fmt.Fprintf(w, "%+v", data.Form)
}

func TestDecoder(t *testing.T) {
	d := midware.Decoder{
		RequestDecoder: formdecoder.New(),
		FieldName:      "Form",
	}
	m := midware.New()
	m.Initilizer(d)
	h := m.Handle(personFunc)

	form := url.Values{
		"Name":  {"Adam"},
		"Age":   {"39"},
		"Admin": {"false"},
	}
	r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, "&{Name:Adam Age:39 Admin:false}", w.Body.String())
}

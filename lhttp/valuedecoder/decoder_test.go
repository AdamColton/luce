package valuedecoder_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamcolton/luce/lhttp/lhttptest"
	"github.com/adamcolton/luce/lhttp/valuedecoder"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name  string
	Age   int
	Admin bool
}

func TestFormError(t *testing.T) {
	fd := valuedecoder.Form()
	p := &Person{}
	r := httptest.NewRequest("POST", "/", nil)
	r.Body = nil
	err := fd.Decode(p, r)
	assert.Error(t, err)
}

func TestDecode(t *testing.T) {
	adam37 := &Person{
		Name:  "Adam",
		Age:   37,
		Admin: true,
	}
	tt := map[string]struct {
		expected *Person
		r        *http.Request
		d        valuedecoder.Decoder
	}{
		"basic-form": {
			expected: adam37,
			r:        lhttptest.NewRequest("/", adam37).POST(),
			d:        valuedecoder.Form(),
		},
		"basic-query": {
			expected: adam37,
			r:        lhttptest.NewRequest("/", adam37).GET(),
			d:        valuedecoder.Query(),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			p := &Person{}
			err := tc.d.Decode(p, tc.r)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, p)
		})
	}
}

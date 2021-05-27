package formdecoder

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name  string
	Age   int
	Admin bool
}

func TestFormDecoder(t *testing.T) {
	tt := map[string]struct {
		expected *Person
		form     map[string][]string
	}{
		"basic": {
			expected: &Person{
				Name:  "Adam",
				Age:   37,
				Admin: true,
			},
			form: map[string][]string{
				"Name":  {"Adam"},
				"Age":   {"37"},
				"Admin": {"true"},
			},
		},
		"admin-false": {
			expected: &Person{
				Name:  "Adam",
				Age:   37,
				Admin: false,
			},
			form: map[string][]string{
				"Name":  {"Adam"},
				"Age":   {"37"},
				"Admin": {"false"},
			},
		},
		"no-age": {
			expected: &Person{
				Name:  "Adam",
				Age:   0,
				Admin: false,
			},
			form: map[string][]string{
				"Name":  {"Adam"},
				"Admin": {"false"},
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			p := &Person{}
			r := httptest.NewRequest("POST", "/", strings.NewReader(url.Values(tc.form).Encode()))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			err := New().Decode(p, r)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, p)
		})
	}
}

func TestFormError(t *testing.T) {
	fd := New()
	p := &Person{}
	r := httptest.NewRequest("POST", "/", nil)
	r.Body = nil
	err := fd.Decode(p, r)
	assert.Error(t, err)
}

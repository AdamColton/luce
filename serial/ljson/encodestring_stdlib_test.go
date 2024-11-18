package ljson_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

// stdlib encodes s with the standard library's json encoder.
func stdlib(s string, escapeHTML bool) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(escapeHTML)
	enc.Encode(s)
	return strings.TrimSuffix(buf.String(), "\n")
}

// EncodeString is taken from the standard library, so it should agree with it.
func TestEncodeStringMatchesStdlib(t *testing.T) {
	inputs := []string{
		"",
		"plain text",
		`quotes " and \ backslash`,
		"control \x00\x01\x1f \n \r \t \b \f",
		"html <b>&</b>",
		"unicode: é世界 😀",
		"separators: \u2028 \u2029 end",
		"del \x7f",
	}
	for _, s := range inputs {
		for _, esc := range []bool{true, false} {
			expected := stdlib(s, esc)
			assert.Equal(t, expected, string(ljson.EncodeString(nil, s, esc)), "%q %v", s, esc)
			assert.Equal(t, expected, string(ljson.EncodeString(nil, []byte(s), esc)), "%q %v", s, esc)
		}
	}
}

// Invalid UTF-8 is written as an escaped replacement character. (Newer versions of the standard
// library write the replacement character itself; both are valid json.)
func TestEncodeStringInvalidUTF8(t *testing.T) {
	for _, esc := range []bool{true, false} {
		assert.Equal(t, `"a\ufffd\ufffdb"`, string(ljson.EncodeString(nil, "a\xff\xfeb", esc)))
		assert.Equal(t, `"a\ufffd\ufffdb"`, string(ljson.EncodeString(nil, []byte("a\xff\xfeb"), esc)))
	}
}

func TestEncodeStringAppends(t *testing.T) {
	got := ljson.EncodeString([]byte("x:"), "a", true)
	assert.Equal(t, `x:"a"`, string(got))
}

package gothicgo

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	"github.com/adamcolton/luce/lerr"
)

// IsExported checks if the first character of name is upper case.
func IsExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

// PrefixWriteToString takes a PrefixWriterTo and Prefixer and uses them to
// generate a string. It is intended for testing rather than production use.
func PrefixWriteToString(w PrefixWriterTo, p Prefixer) string {
	buf := bytes.NewBuffer(nil)
	_, err := w.PrefixWriteTo(buf, p)
	lerr.Panic(err)
	return buf.String()
}

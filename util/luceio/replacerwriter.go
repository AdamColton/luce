package luceio

import (
	"io"
	"strings"
)

// ReplacerWriter invokes a Replacer on any data before sending it to the
// underlying Writer.
type ReplacerWriter struct {
	Writer   io.Writer
	Replacer *strings.Replacer
}

// Write casts the []byte to string, calls Replacer and writes to the underlying
// writer.
func (rw ReplacerWriter) Write(b []byte) (int, error) {
	return rw.Replacer.WriteString(rw.Writer, string(b))
}

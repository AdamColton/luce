package luceio

import (
	"io"
	"strings"
)

// Replacer is fulfilled by strings.Replacer. WriteString writes s to w with all
// replacements performed.
type Replacer interface {
	WriteString(w io.Writer, s string) (int, error)
}

// NewReplacer joins a strings.Replacer with a Writer so that the replacer
// will be applied to all writes.
func NewReplacer(w io.Writer, oldnew ...string) ReplacerWriter {
	return ReplacerWriter{
		Writer:   w,
		Replacer: strings.NewReplacer(oldnew...),
	}
}

// ReplacerWriter invokes a Replacer on any data before sending it to the
// underlying Writer.
type ReplacerWriter struct {
	Writer io.Writer
	Replacer
}

// Write casts the []byte to string, calls Replacer and writes to the underlying
// writer.
func (rw ReplacerWriter) Write(b []byte) (int, error) {
	return rw.Replacer.WriteString(rw.Writer, string(b))
}

// WriteString fulfills
func (rw ReplacerWriter) WriteString(str string) (int, error) {
	return rw.Replacer.WriteString(rw.Writer, str)
}

package luceio

import (
	"io"
	"strings"
)

type Replacer interface {
	WriteString(io.Writer, string) (int, error)
}

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

func (rw ReplacerWriter) WriteString(str string) (int, error) {
	return rw.Write([]byte(str))
}

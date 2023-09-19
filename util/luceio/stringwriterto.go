package luceio

import "io"

// StringWriterTo fulfills io.WriterTo and writes the string to the writer.
type StringWriterTo string

// WriteTo fulfills the WriterTo interface.
func (s StringWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(s))
	return int64(n), err
}

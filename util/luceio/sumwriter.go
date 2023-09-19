package luceio

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// SumWriter is helper that wraps a Writer and sums the bytes written. If it
// encounters an error, it will stop writing.
type SumWriter struct {
	io.Writer
	// Cache holds bytes that will be written before the next write operation.
	// If no write operation is executed, they will never be written.
	Cache []byte
	Sum   int64
	Err   error
}

// NewSumWriter takes a Writer and returns a SumWriter
func NewSumWriter(w io.Writer) *SumWriter {
	return &SumWriter{Writer: w}
}

// BufferSumWriter creates a new buffer passing it into the SumWriter and
// returns both.
func BufferSumWriter() (*bytes.Buffer, *SumWriter) {
	buf := bytes.NewBuffer(nil)
	return buf, NewSumWriter(buf)
}

// WriteString writes a string to underlying Writer
func (s *SumWriter) WriteString(str string) (int, error) {
	return s.Write([]byte(str))
}

// WriteStrings writes strings to underlying Writer
func (s *SumWriter) WriteStrings(strs ...string) (int, error) {
	d := s.Sum
	for _, str := range strs {
		_, err := s.Write([]byte(str))
		if err != nil {
			return s.Diff(d)
		}
	}
	return s.Diff(d)
}

// WriteRune writes a rune to underlying Writer
func (s *SumWriter) WriteRune(r rune) { s.Write([]byte(string(r))) }

// Write fulfills io.Write
func (s *SumWriter) Write(b []byte) (int, error) {
	c := s.FlushCache()
	if s.Err != nil {
		return c, s.Err
	}
	var n int
	n, s.Err = s.Writer.Write(b)
	n += c
	s.Sum += int64(n)
	return n, s.Err
}

// FlushCache causes the Cache to be written and then cleared.
func (s *SumWriter) FlushCache() int {
	c := 0
	if len(s.Cache) > 0 && s.Err == nil {
		c, s.Err = s.Writer.Write(s.Cache)
		s.Cache = s.Cache[:0]
	}
	return c
}

// Rets is a shorthand helper for returns
func (s *SumWriter) Rets() (int64, error) {
	return s.Sum, s.Err
}

// Diff is a shorthand way to handle returns by marking the difference in the
// Sum.
func (s *SumWriter) Diff(d int64) (int, error) {
	return int(s.Sum - d), s.Err
}

// WriteInt uses strconv to write an int
func (s *SumWriter) WriteInt(i int) (int, error) {
	return s.WriteString(strconv.Itoa(i))
}

// WriterTo passes the SumWriter into a WriterTo and captures the character
// length and error.
func (s *SumWriter) WriterTo(w io.WriterTo) (int64, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	var n int64
	n, s.Err = w.WriteTo(s)
	return n, s.Err
}

// Fprint wraps a call to Fprintf.
func (s *SumWriter) Fprint(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(s, format, args...)
}

// Join a list of strings using a provided seperator.
func (s *SumWriter) Join(elems []string, sep string) (int, error) {
	d := s.Sum
	s.WriteString(elems[0])
	for _, e := range elems[1:] {
		s.WriteStrings(sep, e)
	}
	return int(s.Sum - d), s.Err
}

// AppendCacheString will append a string to the current Cache value.
func (s *SumWriter) AppendCacheString(str string) {
	s.AppendCache([]byte(str))
}

// AppendCache will append a byte slice to the current Cache value.
func (s *SumWriter) AppendCache(b []byte) {
	s.Cache = append(s.Cache, b...)
}

// TODO:
func (s *SumWriter) Wrapped() any {
	return s.Writer
}

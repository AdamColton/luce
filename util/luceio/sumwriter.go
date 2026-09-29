package luceio

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/adamcolton/luce/util/liter"
)

// SumWriter is a helper that wraps a Writer and sums the bytes written. If it
// encounters an error, it stores it in Err and stops writing.
type SumWriter struct {
	io.Writer
	// Cache holds bytes that will be written before the next write operation.
	// If no write operation is executed, they will never be written.
	Cache []byte
	// Sum is the total number of bytes written by Write and the methods that
	// call it. It is not updated by FlushCache.
	Sum int64
	// Err is the first error encountered while writing. Once it is set nothing
	// more is written.
	Err error
}

// NewSumWriter takes a Writer and returns a SumWriter.
func NewSumWriter(w io.Writer) *SumWriter {
	return &SumWriter{Writer: w}
}

// BufferSumWriter creates a new buffer passing it into the SumWriter and
// returns both.
func BufferSumWriter() (*bytes.Buffer, *SumWriter) {
	buf := bytes.NewBuffer(nil)
	return buf, NewSumWriter(buf)
}

// WriteString writes a string to the underlying Writer.
func (s *SumWriter) WriteString(str string) (int, error) {
	return s.Write([]byte(str))
}

// WriteStrings writes each of the strings to the underlying Writer, stopping at
// the first error. It returns the number of bytes written by this call.
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

// WriteRune writes a rune to the underlying Writer. Unlike io.RuneWriter, it
// returns nothing; check Err.
func (s *SumWriter) WriteRune(r rune) { s.Write([]byte(string(r))) }

// Write fulfills io.Writer. Anything in Cache is written first and counted in
// the returned n. If Err is already set, nothing is written and Err is
// returned.
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

// FlushCache causes the Cache to be written and then cleared. It returns the
// number of bytes written. It does nothing if Err is set. Sum is not updated;
// Write does that when it flushes the Cache.
func (s *SumWriter) FlushCache() int {
	c := 0
	if len(s.Cache) > 0 && s.Err == nil {
		c, s.Err = s.Writer.Write(s.Cache)
		s.Cache = s.Cache[:0]
	}
	return c
}

// Rets is a shorthand helper for returns. It returns Sum and Err.
func (s *SumWriter) Rets() (int64, error) {
	return s.Sum, s.Err
}

// Diff is a shorthand way to handle returns by marking the difference in the
// Sum. It returns Sum minus d, and Err.
func (s *SumWriter) Diff(d int64) (int, error) {
	return int(s.Sum - d), s.Err
}

// WriteInt uses strconv to write an int in base 10.
func (s *SumWriter) WriteInt(i int) (int, error) {
	return s.WriteString(strconv.Itoa(i))
}

// WriterTo passes the SumWriter into a WriterTo and captures the number of
// bytes written and the error. If Err is already set, it returns 0 and Err
// without calling w.
func (s *SumWriter) WriterTo(w io.WriterTo) (int64, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	var n int64
	n, s.Err = w.WriteTo(s)
	return n, s.Err
}

// Fprint wraps a call to fmt.Fprintf, writing the formatted string.
func (s *SumWriter) Fprint(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(s, format, args...)
}

// Join writes a list of strings using a provided separator. It returns the
// number of bytes written by this call. An empty list writes nothing.
func (s *SumWriter) Join(elems []string, sep string) (int, error) {
	d := s.Sum
	if len(elems) == 0 {
		return 0, s.Err
	}
	s.WriteString(elems[0])
	for _, e := range elems[1:] {
		s.WriteStrings(sep, e)
	}
	return int(s.Sum - d), s.Err
}

// Iter writes the strings from an iterator using sep as a separator. Writing
// starts from the iterator's current value, it is not reset. As with Join, an
// empty iterator writes nothing. It returns the number of bytes written by this
// call.
func (s *SumWriter) Iter(elems liter.Iter[string], sep string) (int, error) {
	d := s.Sum
	cur, done := elems.Cur()
	if !done {
		s.WriteString(cur)
		for cur, done = elems.Next(); !done; cur, done = elems.Next() {
			s.WriteStrings(sep, cur)
		}
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

// Wrapped fulfills upgrade.Wrapper, returning the underlying Writer so it can
// be upgraded to any other interface it fulfills.
func (s *SumWriter) Wrapped() any {
	return s.Writer
}

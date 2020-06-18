package luceio

import "io"

// SumWriter is helper that wraps a Writer and sums the bytes written. If it
// encounters an error, it will stop writing.
type SumWriter struct {
	io.Writer
	Sum int64
	Err error
}

// NewSumWriter takes a Writer and returns a SumWriter
func NewSumWriter(w io.Writer) *SumWriter {
	return &SumWriter{Writer: w}
}

// WriteString writes a string to underlying Writer
func (s *SumWriter) WriteString(str string) (int, error) {
	return s.Write([]byte(str))
}

// WriteStrings writes strings to underlying Writer
func (s *SumWriter) WriteStrings(strs ...string) (int, error) {
	var sum int
	for _, str := range strs {
		i, err := s.Write([]byte(str))
		if err != nil {
			return sum, err
		}
		sum += i
	}
	return sum, nil
}

// WriteRune writes a rune to underlying Writer
func (s *SumWriter) WriteRune(r rune) { s.Write([]byte(string(r))) }

// Write fulfills io.Write
func (s *SumWriter) Write(b []byte) (int, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	var n int
	n, s.Err = s.Writer.Write(b)
	s.Sum += int64(n)
	return n, s.Err
}

// Rets is a shorthand helper for returns
func (s *SumWriter) Rets() (int64, error) {
	return s.Sum, s.Err
}

package lstr

import (
	"unicode/utf8"

	"github.com/adamcolton/luce/util/iter"
)

// Scanner is used it iterate over the runes in a string. Scanner fulfills
// Iter[rune].
type Scanner struct {
	// Str is the string as a byte slice
	Str []byte
	// I is the index of the current rune
	I int
	// Count is the number of runes read
	Count int
	// Rune at the current position
	Rune rune
	// Size of the current rune
	Size int
}

// ScannerFactory creates a function to fulfill iter.Factory[rune] that is
// backed by a Scanner.
func ScannerFactory(str string) iter.Factory[rune] {
	return func() (i iter.Iter[rune], r rune, done bool) {
		i = NewScanner(str)
		r, done = i.Cur()
		return
	}
}

// NewScanner creates a Scanner for the string.
func NewScanner(str string) *Scanner {
	s := &Scanner{
		Str: []byte(str),
	}
	s.read()
	return s
}

func (s *Scanner) read() {
	s.Rune, s.Size = utf8.DecodeRune(s.Str[s.I:])
}

// Next increments the Scanner.
func (s *Scanner) Next() (r rune, done bool) {
	if s.Size > 0 {
		s.Count++
		s.I += s.Size
		s.read()
	}
	return s.Cur()
}

// Cur returns the current rune and a bool indicating if the scanner is done.
func (s *Scanner) Cur() (r rune, done bool) {
	return s.Rune, s.Done()
}

// Idx returns the index of the current rune. The index will skip values if
// there are runes that take multiple bytes.
func (s *Scanner) Idx() int {
	return s.I
}

// Reset the scanner back to the start of the string.
func (s *Scanner) Reset() {
	s.I = 0
	s.Count = 0
	s.read()
}

// Done returns true if I is at or past the end of the string.
func (s *Scanner) Done() bool {
	return s.I >= len(s.Str)
}

// Peek passes the current rune into m and returns the bool.
func (s *Scanner) Peek(m Matcher) bool {
	return m.Matches(s.Rune)
}

// Match checks if the current rune against m and returns the bool. If it is
// is a match, the Scanner moves to the next value.
func (s *Scanner) Match(m Matcher) bool {
	matched := s.Peek(m)
	if matched {
		s.Next()
	}
	return matched
}

// Many checks if the current rune againts m and if it matches iterates until it
// reaches a rune that does not or reaches the end of the string. The returned
// bool indicates if it found at least one match.
func (s *Scanner) Many(m Matcher) bool {
	if !s.Match(m) {
		return false
	}
	for s.Match(m) && !s.Done() {
	}
	return true
}

// Iter returns a wrapped iterator.
func (s *Scanner) Iter() iter.Wrapper[rune] {
	return iter.Wrapper[rune]{s}
}

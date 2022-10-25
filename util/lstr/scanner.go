package lstr

import (
	"unicode/utf8"
)

type Scanner struct {
	Str        []byte
	Idx, Count int
	Rune       rune
	Size       int
}

func NewScanner(str string) *Scanner {
	s := &Scanner{
		Str: []byte(str),
	}
	s.Rune, s.Size = utf8.DecodeRune(s.Str[s.Idx:])
	return s
}

func (s *Scanner) Next() {
	if s.Size > 0 {
		s.Count++
		s.Rune, s.Size = utf8.DecodeRune(s.Str[s.Idx:])
		s.Idx += s.Size
	}
}

func (s *Scanner) Reset() {
	s.Idx = 0
	s.Count = 0
	s.Rune, s.Size = utf8.DecodeRune(s.Str[s.Idx:])
}

func (s *Scanner) Done() bool {
	return s.Idx >= len(s.Str)
}

func (s *Scanner) Peek(c Matcher) bool {
	return c.Matches(s.Rune)
}

func (s *Scanner) Match(c Matcher) bool {
	matched := s.Peek(c)
	if matched {
		s.Next()
	}
	return matched
}

func (s *Scanner) Many(c Matcher) bool {
	if !s.Match(c) {
		return false
	}
	for s.Match(c) {
	}
	return true
}

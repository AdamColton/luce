package txtidx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type scanner struct {
	s              []byte
	idx            int
	r              rune
	size           int
	isLetterNumber bool
}

func newScanner(str string) *scanner {
	s := &scanner{
		s: []byte(str),
	}
	s.next()
	return s
}

func (s *scanner) next() {
	s.idx += s.size
	s.r, s.size = utf8.DecodeRune(s.s[s.idx:])
	s.isLetterNumber = unicode.IsLetter(s.r) || unicode.IsNumber(s.r)
}

func (s *scanner) matchLetterNumber(b bool) {
	for ; !s.done() && s.isLetterNumber != b; s.next() {
	}
}

func (s *scanner) done() bool {
	return s.idx >= len(s.s)
}

func (s *scanner) str(start, end int) string {
	return string(s.s[start:end])
}

type search struct {
	words, exact []string
}

func (s *scanner) buildSearch() search {
	out := search{}
	exactStart := -1
	start := -1
	for ; !s.done(); s.next() {
		if s.r == '"' {
			if exactStart == -1 {
				exactStart = s.idx + 1
			} else {
				out.exact = append(out.exact, string(s.s[exactStart:s.idx]))
				exactStart = -1
			}
		}
		if s.isLetterNumber && start == -1 {
			start = s.idx
		} else if !s.isLetterNumber && start != -1 {
			out.words = append(out.words, strings.ToLower(string(s.s[start:s.idx])))
			start = -1
		}
	}
	if start != -1 {
		out.words = append(out.words, strings.ToLower(string(s.s[start:s.idx])))
	}
	return out
}

func (s *scanner) parse(str string) ([]byte, []string) {
	s.matchLetterNumber(true)
	start := s.s[:s.idx]

	var words []string
	for !s.done() {
		start := s.idx
		s.matchLetterNumber(false)
		s.matchLetterNumber(true)
		words = append(words, s.str(start, s.idx))
	}
	return start, words
}

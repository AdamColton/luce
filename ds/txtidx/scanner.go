package txtidx

import (
	"strings"

	"github.com/adamcolton/luce/util/lstr"
)

type search struct {
	words, exact []string
}

var (
	mIsLetterNumber    = lstr.Or{lstr.IsLetter, lstr.IsNumber}
	mNotIsLetterNumber = lstr.Not{mIsLetterNumber}
)

func buildSearch(s *lstr.Scanner) search {
	out := search{}
	exactStart := -1
	start := -1
	for done := false; !done; _, done = s.Next() {
		if s.Rune == '"' {
			if exactStart == -1 {
				exactStart = s.I + 1
			} else {
				out.exact = append(out.exact, string(s.Str[exactStart:s.I]))
				exactStart = -1
			}
		}
		isLetterNumber := s.Peek(mIsLetterNumber)
		if isLetterNumber && start == -1 {
			start = s.I
		} else if !isLetterNumber && start != -1 {
			out.words = append(out.words, strings.ToLower(string(s.Str[start:s.I])))
			start = -1
		}
	}
	if start != -1 {
		out.words = append(out.words, strings.ToLower(string(s.Str[start:s.I])))
	}
	return out
}

func parse(s *lstr.Scanner) ([]byte, []string) {
	s.Many(mNotIsLetterNumber)
	start := s.Str[:s.I]

	var words []string
	for !s.Done() {
		start := s.I
		s.Many(mIsLetterNumber)
		s.Many(mNotIsLetterNumber)
		words = append(words, string(s.Str[start:s.I]))
	}
	return start, words
}

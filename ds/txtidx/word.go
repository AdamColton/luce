package txtidx

import "strings"

type word struct {
	str string
	wordIDX
}

type wordIDX uint32

// str must start with letterNumber but can have trailing non-letter number
func root(str string) string {
	s := newScanner(str)
	s.matchLetterNumber(false)
	return strings.ToLower(s.str(0, s.idx))
}

package document

import (
	"strings"

	"github.com/adamcolton/luce/util/lstr"
)

var (
	mIsLetterNumber = lstr.Or{lstr.IsLetter, lstr.IsNumber}
)

// Root finds the prefix of the string containing letters and numbers, in lower
// case.
func Root(str string) string {
	s := lstr.NewScanner(str)
	s.Many(mIsLetterNumber)
	return strings.ToLower(string(s.Str[0:s.I]))
}

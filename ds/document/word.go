package document

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/util/lstr"
)

var (
	mIsLetterNumber = lstr.Or{lstr.IsLetter, lstr.IsNumber}
)

// Root find the prefix of the string containing letters and numbers.
func Root(str string) string {
	s := lstr.NewScanner(str)
	s.Many(mIsLetterNumber)
	return strings.ToLower(string(s.Str[0:s.I]))
}

// RootVariant find the prefix of the string containing letters and numbers and
// the Variant to convert the root back to the original input.
func RootVariant(str string) (string, Variant) {
	root := Root(str)
	return root, findVariant(root, str)
}

// Variant encodes the casing of a word and the non-alphanumeric characters
// that follow the word.
type Variant []byte

func findVariant(root, str string) Variant {
	rs := []rune(root)
	b := []byte(str)
	suffix := str[len(root):]

	caseData := &rye.Bits{
		Data: make([]byte, 0, divUp(len(rs), 8)+len(suffix)),
	}
	for _, rr := range rs {
		r, ln := utf8.DecodeRune(b)
		b = b[ln:]
		if r != rr {
			caseData.Write(1)
		} else {
			caseData.Write(0)
		}
	}

	return append(caseData.Data, suffix...)
}

// Apply a variant to a word. It is expected that the root is all lower case.
// The casing will be changed according the variant and non-alphanumeric
// runes will be appended.
func (v Variant) Apply(root string, buf []byte) string {
	b := &rye.Bits{
		Data: v,
	}
	in := []byte(root)
	for len(in) > 0 {
		r, size := utf8.DecodeRune(in)
		in = in[size:]
		if b.Read() == 1 {
			r = unicode.ToUpper(r)
		}

		buf = append(buf, string(r)...)

	}

	buf = append(buf, v[divUp(b.Idx, 8):]...)
	return string(buf)
}

// divUp division round up
func divUp(a, b int) int {
	out := a / b
	if out*b != a {
		out++
	}
	return out
}

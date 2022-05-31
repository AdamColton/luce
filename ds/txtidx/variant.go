package txtidx

import (
	"unicode"
	"unicode/utf8"

	"github.com/adamcolton/luce/serial/rye"
)

type varIDX uint32

type variant []byte

func findVariant(root, str string) variant {
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

// divUp division round up
func divUp(a, b int) int {
	out := a / b
	if out*b != a {
		out++
	}
	return out
}

// todo: this should take a buffer
func (v variant) apply(rt string) string {
	b := &rye.Bits{
		Data: v,
	}
	in := []byte(rt)
	var out []byte
	for len(in) > 0 {
		r, size := utf8.DecodeRune(in)
		in = in[size:]
		if b.Read() == 1 {
			r = unicode.ToUpper(r)
		}

		out = append(out, string(r)...)

	}

	out = append(out, v[divUp(b.Idx, 8):]...)
	return string(out)
}

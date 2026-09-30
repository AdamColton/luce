package lstr

import (
	"strings"

	"github.com/adamcolton/luce/ds/slice"
)

// Seperator is used for string operations with a separator.
type Seperator string

// JoinLen returns the length of joining the elements with a single Seperator.
// This is used by BufJoin to allocate the correct size slice for the output.
func (s Seperator) JoinLen(elems []string) int {
	if len(elems) == 0 {
		return 0
	}
	sep := string(s)
	sln := len(s)
	prev := elems[0]
	pln := len(prev)
	ln := pln
	for _, e := range elems[1:] {
		if e == "" {
			continue
		}
		ps := len(prev) >= sln && prev[pln-sln:] == sep
		prev, pln = e, len(e)
		es := pln >= sln && e[:sln] == sep
		ln += pln
		if ps && es {
			ln -= sln
		} else if !ps && !es {
			ln += sln
		}
	}
	return ln
}

// BufJoin joins elems making sure there is a single Seperator between each elem
// and will use buf if it has adequate capacity. If an elem ends with the
// Seperator and the next starts with it, only one is kept and if neither has it,
// one is added. Empty elems are skipped.
func (s Seperator) BufJoin(elems []string, buf []byte) string {
	if len(elems) == 0 {
		return ""
	}

	sep := string(s)
	sln := len(s)
	ln := s.JoinLen(elems)
	out := slice.NewBuffer(buf).Empty(ln)

	out = append(out, elems[0]...)
	for _, e := range elems[1:] {
		if e == "" {
			continue
		}
		oln := len(out)
		ps := oln >= sln && string(out[oln-sln:]) == sep
		es := len(e) >= sln && e[:sln] == sep
		if ps && es {
			out = append(out, e[sln:]...)
		} else {
			if !ps && !es {
				out = append(out, sep...)
			}
			out = append(out, e...)
		}
	}
	return string(out)
}

// Join elems making sure there is a single Seperator between each elem. See
// BufJoin.
func (s Seperator) Join(elems ...string) string {
	return s.BufJoin(elems, nil)
}

// Index is a wrapper around strings.Index, returning the index of the first
// Seperator in str or -1.
func (s Seperator) Index(str string) int {
	return strings.Index(str, string(s))
}

// Split is a wrapper around strings.Split, splitting str at every Seperator.
func (s Seperator) Split(str string) slice.Slice[string] {
	return strings.Split(str, string(s))
}

// Trailing makes sure str ends with the Seperator if want is true and does not
// end with it if want is false. If str already is as wanted it is returned
// unchanged. Only one Seperator is added or removed, so with a "/" Seperator
// "a//" becomes "a/". An empty str is returned as it is, there is nothing for
// the Seperator to trail.
func (s Seperator) Trailing(str string, want bool) string {
	if str == "" {
		return str
	}
	sep := string(s)
	has := strings.HasSuffix(str, sep)
	if want && !has {
		return str + sep
	}
	if !want && has {
		return str[:len(str)-len(sep)]
	}
	return str
}

const (
	// NewLine seperator
	NewLine Seperator = "\n"
	// Space seperator
	Space Seperator = " "
)

// Strings applies the Seperator to split str and returns an instance of
// Strings.
func (s Seperator) Strings(str string) *Strings {
	return NewStrings(strings.Split(str, string(s)))
}

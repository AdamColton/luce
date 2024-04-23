package lstr

import (
	"strings"

	"github.com/adamcolton/luce/ds/slice"
)

// Seperator is used for string operations with a seperator
type Seperator string

// JoinLen returns the length of joining the elements with a single Seperator.
// This is used by Join to allocate the correct size slice for the output.
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
// and will use buf if it has adequate capacity.
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
			out = append(out, e[1:]...)
		} else {
			if !ps && !es {
				out = append(out, sep...)
			}
			out = append(out, e...)
		}
	}
	return string(out)
}

// Join elems making sure there is a single Seperator between each elem.
func (s Seperator) Join(elems ...string) string {
	return s.BufJoin(elems, nil)
}

// Index is a wrapper around strings.Index
func (s Seperator) Index(str string) int {
	return strings.Index(str, string(s))
}

// Split is a wrapper around strings.Split
func (s Seperator) Split(str string) slice.Slice[string] {
	return strings.Split(str, string(s))
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

// Joiner wraps strings.Join. This allows it to be invoked lazily.
type Joiner struct {
	Elems []string
	Seperator
}

// Joiner creates an instance of Joiner with the given elements.
func (s Seperator) Joiner(elems ...string) *Joiner {
	return &Joiner{
		Elems:     elems,
		Seperator: s,
	}
}

// Stirng fulfills fmt.Stringer and calls strings.Join.
func (j *Joiner) String() string {
	return strings.Join(j.Elems, string(j.Seperator))
}

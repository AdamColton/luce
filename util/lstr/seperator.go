package lstr

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

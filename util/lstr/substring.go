package lstr

// SubStrings represents a slice of substrings by index.
type SubStrings [][2]uint

// Slice applies the SubStrings to str and returns a slice of strings.
func (s SubStrings) Slice(str string) []string {
	out := make([]string, len(s))
	for i, idx := range s {
		out[i] = str[idx[0]:idx[1]]
	}
	return out
}

// SubStringBySplit creates a SubStrings slice by splitting the string into
// sections so that the end of each range is equal to the start of the next.
func SubStringBySplit(splits []int) SubStrings {
	ln := len(splits)
	out := make(SubStrings, ln)
	ln--
	for i, s := range splits[:ln] {
		out[i][1] = uint(s)
		out[i+1][0] = uint(s)
	}
	out[ln][1] = uint(splits[ln])
	return out
}

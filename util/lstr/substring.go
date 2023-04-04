package lstr

// SubStrings represents a list of substrings of a string, each as a byte range
// from the first index up to but not including the second.
type SubStrings [][2]uint

// Slice applies the SubStrings to str and returns a slice of strings. It panics
// if a range is outside of str.
func (s SubStrings) Slice(str string) []string {
	out := make([]string, len(s))
	for i, idx := range s {
		out[i] = str[idx[0]:idx[1]]
	}
	return out
}

// SubStringBySplit creates a SubStrings slice by splitting the string into
// sections so that the end of each range is equal to the start of the next. The
// first section starts at 0 and each split is the end of a section, so the last
// split is the end of the last section. It panics if splits is empty.
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

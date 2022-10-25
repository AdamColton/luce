package lstr

type SubStrings [][2]uint

func (s SubStrings) Slice(str string) []string {
	out := make([]string, len(s))
	for i, idx := range s {
		out[i] = str[idx[0]:idx[1]]
	}
	return out
}

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

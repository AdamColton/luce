package lstr

import "strings"

func Len(s string) int {
	return len(s)
}

func NewRemover(rm ...string) *strings.Replacer {
	s := make([]string, len(rm)*2)
	for i, r := range rm {
		s[i*2] = r
		s[i*2+1] = ""
	}
	return strings.NewReplacer(s...)
}

package luceio

import "strings"

// Join is syntact sugar for calling strings.Join. The last argument is still
// used as the seperator, but it doesn't require secifying the slice.
func Join(strs ...string) string {
	ln := len(strs)
	if ln < 2 {
		return ""
	}
	if ln == 2 {
		return strs[0]
	}
	return strings.Join(strs[:ln-1], strs[ln-1])
}

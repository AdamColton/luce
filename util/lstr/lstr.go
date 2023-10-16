package lstr

// Len is a strongly typed version of the builtin len for strings. It is useful
// when a func value is needed.
func Len(s string) int {
	return len(s)
}

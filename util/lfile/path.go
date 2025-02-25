package lfile

import "github.com/adamcolton/luce/util/lstr"

// Name returns the last portion of a path as its name.
//   - "/foo/bar.txt" => "/foo/","bar.txt"
//   - "/foo/bar/" => "/foo/","bar"
//   - "foo.txt" => "", "foo.txt"
//   - foo/ => "", "foo"
//
// The second returned value is the name and the first
// is the preceding portion.
func Name(path string) (string, string) {
	end := len(path) - 1
	if end < 0 {
		return "", ""
	}
	for end > 0 && path[end] == '/' {
		end--
	}
	start := end - 1
	for ; start >= 0 && path[start] != '/'; start-- {
	}
	start++
	return path[:start], path[start : end+1]
}

// Slash makes sure path ends with a slash if trail is true and does not if it
// is false. Only one slash is added or removed, and an empty path stays empty.
// It is lstr.Seperator("/").Trailing.
func Slash(path string, trail bool) string {
	return lstr.Seperator("/").Trailing(path, trail)
}

package lfile

import (
	"path/filepath"
	"strings"

	"github.com/adamcolton/luce/util/lstr"
)

// PathLength is used to trim filenames to a set number of parts.
type PathLength int

var separator = string(filepath.Separator)

// Trim the filename so it will have at most PathLength number of parts,
// including the filename. The returned value will never begin with
// filepath.Separator. Filenames are real files of the operating system, so
// unlike the paths that go through a CoreFS, the parts are separated by
// filepath.Separator rather than a slash.
func (pln PathLength) Trim(filename string) string {
	if pln <= 0 {
		return ""
	}
	idx := len(filename)
	for c := pln - 1; c >= 0; c-- {
		idx = strings.LastIndex(filename[:idx], separator)
		if idx < 0 {
			break
		}
	}
	if idx+1 < len(filename) && filename[idx+1] == filepath.Separator {
		idx++
	}

	return filename[idx+1:]
}

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

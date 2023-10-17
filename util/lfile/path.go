package lfile

import (
	"path/filepath"
	"strings"
)

// PathLength is used to trim filenames to a set number of parts.
type PathLength int

var separator = string(filepath.Separator)

// Trim the filname so it will have at most PathLength number of parts,
// including the filename. The returned value will never begin with
// filepath.Separator.
func (pln PathLength) Trim(filename string) string {
	idx := len(filename)
	for c := pln - 1; c >= 0; c-- {
		idx = strings.LastIndex(filename[:idx], separator)
		if idx < 0 {
			break
		}
	}
	if filename[idx+1] == filepath.Separator {
		idx++
	}

	return filename[idx+1:]
}

// Name returns the last portion of a path as it's name.
//   - "/foo/bar.txt" => "/foo/","bar.txt"
//   - "/foo/bar/" => "/foo/","bar"
//   - "foo.txt" => "", "foo.txt"
//   - foo/ => "", "foo"
//
// The second returned value is the name and the first
// is the preceeding portion.
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

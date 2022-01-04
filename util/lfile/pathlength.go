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

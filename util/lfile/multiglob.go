// Package lfiles provides utility functions around file and directory
// operations.
package lfile

import (
	"path/filepath"
)

// Glob is a reference to filepath.Glob. It is left exposed for testing.
var Glob = filepath.Glob

// MultiGlob performs multiple filepath.Glob operation and merges the values
// into a single slice with no duplicates.
type MultiGlob []string

// Paths found using MultiGlob.
func (mg MultiGlob) Paths() ([]string, error) {
	unique := make(map[string]bool)

	for _, p := range mg {
		files, err := Glob(p)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			unique[f] = true
		}
	}

	out := make([]string, 0, len(unique))
	for f := range unique {
		out = append(out, f)
	}
	return out, nil
}

// Iter will iterate over the files found by the MultiGlob.
func (mg MultiGlob) Iterator() (Iterator, bool) {
	fs, err := mg.Paths()
	if err != nil {
		return &pathsIterator{err: err}, true
	}

	return Paths(fs).Iterator()
}

// Recursive adds two glob patterns, one to the base and one recursive.
func (mg MultiGlob) Recursive(base, pattern string) MultiGlob {
	return append(mg,
		base+pattern,
		base+"**/"+pattern,
	)
}

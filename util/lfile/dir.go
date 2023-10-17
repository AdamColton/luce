package lfile

import (
	"sort"
)

// DirContents is returned from GetDirContents.
type DirContents struct {
	Name, Path string
	SubDirs    []string
	Files      []string
}

// GetDirContents splits the full path of the Dir into it's name and partial
// path and returns the sub-dirs and files in sorted slices.
func GetDirContents(f Dir) (*DirContents, error) {
	fs, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	fIdx := len(fs)
	out := make([]string, fIdx)
	dIdx := 0
	fIdx--
	for _, f := range fs {
		if f.IsDir() {
			out[dIdx] = f.Name()
			dIdx++
		} else {
			out[fIdx] = f.Name()
			fIdx--
		}
	}
	sort.Strings(out[:dIdx])
	sort.Strings(out[dIdx:])

	dc := &DirContents{
		SubDirs: out[:dIdx],
		Files:   out[dIdx:],
	}
	dc.Path, dc.Name = Name(f.Name())
	return dc, nil
}

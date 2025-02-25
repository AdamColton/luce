package lfile

import (
	"io"
	"io/fs"
	"sort"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/upgrade"
)

// CoreFsStat describes a file. It uses the Stat of the file system if it has one,
// otherwise it opens the file, describes it and closes it.
func CoreFsStat(cf CoreFS, name string) (fs.FileInfo, error) {
	if s, ok := upgrade.To[FSStater](cf); ok {
		return s.Stat(name)
	}
	return fs.Stat(cf, name)
}

var getNames = morph.NewValAll(fs.DirEntry.Name)

// ReadDirNames returns the sorted names of the entries of a directory. If the
// opened directory can list names (a Dir or DirNameReader) that is used,
// otherwise the directory is read with fs.ReadDir. It is an error if the path
// is not a directory.
func ReadDirNames(r FSOpener, dirname string) ([]string, error) {
	f, err := r.Open(dirname)
	if err != nil {
		return nil, err
	}
	nr, ok := upgrade.To[DirNameReader](f)
	if !ok {
		f.Close()
		de, err := fs.ReadDir(r, dirname)
		if err != nil {
			return nil, err
		}
		return getNames.Slice(de, nil), nil
	}

	names, err := nr.Readdirnames(-1)
	f.Close()
	if err != nil && !lerr.Except(err, io.EOF) {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

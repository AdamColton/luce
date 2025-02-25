package lfile

import (
	"io"
	"io/fs"
	"sort"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/upgrade"
)

func CoreFsStat(cf CoreFS, name string) (fs.FileInfo, error) {
	if s, ok := upgrade.To[FSStater](cf); ok {
		return s.Stat(name)
	}

	f, err := cf.Open(name)
	if err != nil {
		return nil, err
	}
	out, err := f.Stat()
	f.Close()
	return out, err
}

var getNames = slice.ForAll(fs.DirEntry.Name)

func ReadDirNames(r FSOpener, dirname string) ([]string, error) {
	f, err := r.Open(dirname)
	if err != nil {
		return nil, err
	}
	var names []string
	if drnrdr, ok := upgrade.To[DirNameReader](f); ok {
		names, err = drnrdr.Readdirnames(-1)
	} else {
		s, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if s.IsDir() {
			if fdr, ok := upgrade.To[FSDirReader](r); ok {
				de, err := fdr.ReadDir(dirname)
				if err != nil {
					return nil, err
				}
				names = getNames.Slice(de, nil)
			}
		}
	}
	f.Close()

	if err != nil && !lerr.Except(err, io.EOF) {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

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

func FilesByExt(root string, cfs CoreFS, exts ...string) (slice.Slice[string], error) {
	re := Exts(false, exts...)
	return FilesByRegex(root, cfs, re)
}

func FilesByRegex(root string, cfs CoreFS, re string) (slice.Slice[string], error) {
	m, err := RegexMatch(re, "", "")
	if err != nil {
		return nil, err
	}
	return RootGetFiles(root, cfs, m)
}

func RootGetFiles(root string, cfs CoreFS, m Match) (slice.Slice[string], error) {
	mr := m.Root(root)
	if cfs != nil {
		mr.CoreFS = cfs
	}
	root = Slash(root, true)
	files := GetFiles(root, nil)
	err := RunHandlerSource(mr, files)
	return files.Matches, err
}

func MatchExt(recursive bool, exts ...string) Match {
	re := Exts(false, exts...)
	dirRe := ""
	if !recursive {
		dirRe = ".*"
	}
	return lerr.Must(RegexMatch(re, "", dirRe))
}

package lfile

import (
	"io"
	"io/fs"
	"sort"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/ds/slice"
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

// FilesByExt returns the paths of the files under root that end in one of the
// extensions, given without the dot, and are not in a hidden directory. See
// RootGetFiles.
func FilesByExt(root string, cfs CoreFS, exts ...string) (slice.Slice[string], error) {
	re := Exts(false, exts...)
	return FilesByRegex(root, cfs, re)
}

// FilesByRegex returns the paths of the files under root that match the regular
// expression, or the error if it does not compile. See RootGetFiles.
func FilesByRegex(root string, cfs CoreFS, re string) (slice.Slice[string], error) {
	m, err := RegexMatch(re, "", "")
	if err != nil {
		return nil, err
	}
	return RootGetFiles(root, cfs, m)
}

// RootGetFiles returns the paths of the files under root that m finds, in the
// order MatchRoot walks them. The paths are relative to root, so they do not
// start with it. cfs is the file system to search, and nil is the operating
// system. It also returns any error that ended the walk, along with the paths
// found before it.
func RootGetFiles(root string, cfs CoreFS, m Match) (slice.Slice[string], error) {
	mr := m.Root(root)
	if cfs != nil {
		mr.CoreFS = cfs
	}
	// the paths start with the root as MatchRoot cleaned it, which is nothing
	// for the current directory
	cut := ""
	if mr.Root != "" && mr.Root != "." {
		cut = Slash(mr.Root, true)
	}
	files := GetFiles(cut, nil)
	err := RunHandlerSource(mr, files)
	return files.Matches, err
}

// MatchExt makes a Match for the files that end in one of the extensions, given
// without the dot, that are not in a hidden directory. If recursive is false
// only the files in the root are found, otherwise all the directories below it
// are searched.
func MatchExt(recursive bool, exts ...string) Match {
	re := Exts(false, exts...)
	dirRe := ""
	if !recursive {
		dirRe = ".*"
	}
	return lerr.Must(RegexMatch(re, "", dirRe))
}

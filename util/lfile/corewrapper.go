package lfile

import (
	"io"
	"io/fs"

	"github.com/adamcolton/luce/util/upgrade"
)

// coreWrapper adds what a CoreFS lacks to make it an FSReader.
type coreWrapper struct {
	CoreFS
}

// Wrapped returns the CoreFS that is wrapped so upgrade.To can look through the
// wrapper.
func (ew coreWrapper) Wrapped() any {
	return ew.CoreFS
}

// Stat describes the file. It uses the Stat of the wrapped file system if it
// has one, otherwise it opens the file, describes it and closes it.
func (ew coreWrapper) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(ew.CoreFS, name)
}

// Open opens the file in the wrapped file system and wraps it so it fulfills
// Dir.
func (ew coreWrapper) Open(name string) (fs.File, error) {
	f, err := ew.CoreFS.Open(name)
	if err != nil {
		return nil, err
	}
	out := &fsFileWrapper{
		File: f,
		name: name,
		ew:   ew,
	}
	return out, err
}

// WrapCoreFS makes a CoreFS, such as an embed.FS, into an FSReader. It adds
// Stat to the file system. The files it opens fulfill Dir: they have a Name and
// list their entries through the ReadDir of the file system. Listing a file that
// is not a directory is an error.
func WrapCoreFS(cfs CoreFS) FSReader {
	return coreWrapper{cfs}
}

// fsFileWrapper is a file opened by a coreWrapper. It fulfills Dir.
type fsFileWrapper struct {
	fs.File
	name string
	ew   coreWrapper

	entries []fs.DirEntry
	loaded  bool
	pos     int
}

// Name is the path the file was opened with.
func (f *fsFileWrapper) Name() string {
	return f.name
}

// ReadDir lists the directory as os.File.ReadDir does. If n <= 0 it returns all
// the entries that have not been returned yet. If n > 0 it returns up to n
// entries and io.EOF when there are none left.
func (f *fsFileWrapper) ReadDir(n int) ([]fs.DirEntry, error) {
	if !f.loaded {
		entries, err := f.ew.CoreFS.ReadDir(f.name)
		if err != nil {
			return nil, err
		}
		f.entries = entries
		f.loaded = true
	}

	rest := f.entries[f.pos:]
	if n <= 0 {
		f.pos = len(f.entries)
		return rest, nil
	}
	if len(rest) == 0 {
		return nil, io.EOF
	}
	if n > len(rest) {
		n = len(rest)
	}
	f.pos += n
	return rest[:n], nil
}

// Readdirnames lists the names of the entries as os.File.Readdirnames does. It
// reads the same entries ReadDir would.
func (f *fsFileWrapper) Readdirnames(n int) (names []string, err error) {
	des, err := f.ReadDir(n)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(des))
	for i, de := range des {
		out[i] = de.Name()
	}
	return out, nil
}

// fsCore makes a bare fs.FS a CoreFS by reading through fs.ReadFile and
// fs.ReadDir, which use the file system's own methods when it has them.
type fsCore struct {
	fs.FS
}

// Wrapped returns the fs.FS so upgrade.To can look through the wrapper.
func (c fsCore) Wrapped() any {
	return c.FS
}

func (c fsCore) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(c.FS, name)
}

func (c fsCore) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(c.FS, name)
}

func (c fsCore) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(c.FS, name)
}

// WrapFS makes an fs.FS that has nothing but Open, such as an archive/zip
// reader, into an FSReader. Reading a file, listing a directory and Stat use the
// methods of the file system when it has them and fall back to Open when it
// does not. If the file system is already an FSReader it is returned as it is.
func WrapFS(fsys fs.FS) FSReader {
	if r, ok := upgrade.To[FSReader](fsys); ok {
		return r
	}
	return WrapCoreFS(fsCore{fsys})
}

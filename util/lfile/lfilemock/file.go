package lfilemock

import (
	"io"
	"io/fs"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/adamcolton/luce/ds/lbuf"
	"github.com/adamcolton/luce/lerr"
)

// File mock allows for *os.File to be simulated including various errors. A
// File is opened from a node of the tree by Repository.Open, or directly with
// Node.File or New.
type File struct {
	// FileName is what Name returns. Repository.Open sets it to the path it was
	// given, as os.Open does.
	FileName string
	// Buffer holds the contents of a file. It is nil for a directory.
	*lbuf.Buffer
	// Dir is true for a directory.
	Dir bool
	// DirEntries is what a directory lists.
	DirEntries []os.DirEntry
	os.FileInfo
	// FileSize is the size Stat reports when Buffer is nil.
	FileSize int64
	os.FileMode
	// Mod is the modification time Stat reports.
	Mod time.Time
	// SysData is what Sys returns from the FileInfo Stat creates.
	SysData any
	// Err is returned by every operation on the File.
	Err    error
	dirPos int
}

const (
	// ErrNewType is the panic when New is given contents that are not a string
	// or []byte.
	ErrNewType = lerr.Str("lfilemock.New contents must be string or []byte")
	// ErrParseType is the panic when Parse is given a value that is not a string,
	// []byte, []string or map[string]any.
	ErrParseType = lerr.Str("lfilemock.Parse values must be string, []byte, []string or map[string]any")
	// ErrParseName is the panic when a file or directory is given a name that is
	// empty, "." or ".." or contains a "/".
	ErrParseName = lerr.Str("lfilemock names must not be empty, '.' or '..', or contain '/'")
)

// New creates a file using either a string or []byte. It panics with ErrNewType
// for anything else.
func New(name string, contents any) *File {
	f := (&ByteFile{
		Name: name,
	})
	switch c := contents.(type) {
	case []byte:
		f.Data = lbuf.New(c)
	case string:
		f.Data = lbuf.String(c)
	default:
		panic(ErrNewType)
	}
	return f.File()
}

func (f *File) pathErr(op string, err error) error {
	return &fs.PathError{Op: op, Path: f.FileName, Err: err}
}

// Close returns f.Err.
func (f *File) Close() error {
	return f.Err
}

// Read returns f.Err if that is set, otherwise it wraps f.Buffer.Read. Reading
// a directory is an error.
func (f *File) Read(b []byte) (int, error) {
	if f.Err != nil {
		return 0, f.Err
	}
	if f.Dir {
		return 0, f.pathErr("read", syscall.EISDIR)
	}
	return f.Buffer.Read(b)
}

// Write returns f.Err if that is set, otherwise it wraps f.Buffer.Write. Writing
// a directory is an error.
func (f *File) Write(b []byte) (int, error) {
	if f.Err != nil {
		return 0, f.Err
	}
	if f.Dir {
		return 0, f.pathErr("write", syscall.EBADF)
	}
	return f.Buffer.Write(b)
}

// ReadDir lists the directory as os.File.ReadDir does. If n <= 0 it returns all
// the entries that have not been returned yet. If n > 0 it returns up to n
// entries and io.EOF when there are none left. The position is shared with
// Readdirnames. It returns f.Err if that is set, and an error for a file.
func (f *File) ReadDir(n int) ([]os.DirEntry, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	if !f.Dir {
		return nil, f.pathErr("readdirent", syscall.ENOTDIR)
	}

	rest := f.DirEntries[f.dirPos:]
	if n <= 0 {
		f.dirPos = len(f.DirEntries)
		return rest, nil
	}
	if len(rest) == 0 {
		return nil, io.EOF
	}
	if n > len(rest) {
		n = len(rest)
	}
	f.dirPos += n
	return rest[:n], nil
}

// Name returns FileName
func (f *File) Name() string {
	return f.FileName
}

// Stat creates an instance of lfilemock.FileInfo derived from the instance of
// File. Its name is the last part of FileName and the size of a file is the
// length of its data, which changes as it is written. It returns f.Err if that
// is set.
func (f *File) Stat() (os.FileInfo, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	size := f.FileSize
	if f.Buffer != nil {
		size = int64(f.Buffer.Len())
	}
	return &FileInfo{
		FileName: path.Base(f.FileName),
		FileSize: size,
		FileMode: f.FileMode,
		Mod:      f.Mod,
		Dir:      f.Dir,
		SysData:  f.SysData,
	}, nil
}

// Readdirnames lists the names of the entries as os.File.Readdirnames does. It
// reads the same entries ReadDir would.
func (f *File) Readdirnames(n int) (names []string, err error) {
	des, err := f.ReadDir(n)
	if err != nil {
		return nil, err
	}
	names = make([]string, len(des))
	for i, de := range des {
		names[i] = de.Name()
	}
	return names, nil
}

// DirEntry fulfills os.DirEntry.
type DirEntry struct {
	EntryName string
	Dir       bool
	// FileMode is the type bits Type returns.
	fs.FileMode
	Err error
	fs.FileInfo
}

// Name returns EntryName
func (de *DirEntry) Name() string {
	return de.EntryName
}

// IsDir returns Dir
func (de *DirEntry) IsDir() bool {
	return de.Dir
}

// Type returns FileMode
func (de *DirEntry) Type() fs.FileMode {
	return de.FileMode
}

// Info returns FileInfo and Err
func (de *DirEntry) Info() (fs.FileInfo, error) {
	return de.FileInfo, de.Err
}

// FileInfo fulfills os.FileInfo
type FileInfo struct {
	FileName string
	FileSize int64
	os.FileMode
	Mod time.Time
	Dir bool
	// SysData is what Sys returns.
	SysData any
}

// Name returns Filename
func (fi *FileInfo) Name() string {
	return fi.FileName
}

// Size returns Filesize
func (fi *FileInfo) Size() int64 {
	return fi.FileSize
}

// Mode returns FileMode
func (fi *FileInfo) Mode() os.FileMode {
	return fi.FileMode
}

// ModTime returns Mod
func (fi *FileInfo) ModTime() time.Time {
	return fi.Mod
}

// IsDir returns Dir
func (fi *FileInfo) IsDir() bool {
	return fi.Dir
}

// Sys returns SysData, which is nil unless it was set.
func (fi *FileInfo) Sys() any {
	return fi.SysData
}

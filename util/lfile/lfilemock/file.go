// Package lfilemock provides some mock file objects. This is not meant to be
// exhaustive, it is used to cover the cases needed to test luce. Over time it
// will expand.
package lfilemock

import (
	"io"
	"io/fs"
	"os"
	"syscall"
	"time"

	"github.com/adamcolton/luce/ds/lbuf"
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/lerr"
)

// File mock allows for *os.File to be simulated including various errors.
type File struct {
	FileName string
	*lbuf.Buffer
	Dir        bool
	DirEntries []os.DirEntry
	os.FileInfo
	FileSize int64
	os.FileMode
	Err          error
	readdirnames int
}

const (
	ErrNewType   = lerr.Str("lfilemock.New contents must be string or []byte")
	ErrParseType = lerr.Str("lfilemock.Parse values must be string, []byte or map[string]any")
)

// New creates a file using either a string or []byte.
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

// Close returns f.Err.
func (f *File) Close() error {
	return f.Err
}

// Read returns f.Err is that is set, otherwise it wraps f.Buffer.Read.
func (f *File) Read(b []byte) (int, error) {
	if f.Err != nil {
		return 0, f.Err
	}
	return f.Buffer.Read(b)
}

// ReadDir returns DirEntries and Err.
func (f *File) ReadDir(n int) ([]os.DirEntry, error) {
	return f.DirEntries, f.Err
}

// Name returns Filename
func (f *File) Name() string {
	return f.FileName
}

// Stat creates an instance of lfilemock.FileInfo derived from the instance of
// File.
func (f *File) Stat() (os.FileInfo, error) {
	fi := &FileInfo{
		FileName: f.FileName,
		FileSize: f.FileSize,
		FileMode: f.FileMode,
		Mod:      time.Time{},
		Dir:      f.Dir,
	}
	return fi, f.Err
}

func (f *File) Readdirnames(n int) (names []string, err error) {
	if !f.Dir {
		return nil, &fs.PathError{Op: "readdirent", Path: f.FileName, Err: syscall.ENOTDIR}
	}

	ln := len(f.DirEntries)
	start := f.readdirnames
	if start >= ln {
		return nil, io.EOF
	}

	end := n
	if end <= 0 || end > ln {
		end = ln
	}

	out := list.TransformSlice(f.DirEntries[start:end], os.DirEntry.Name).ToSlice(nil)
	f.readdirnames = end

	return out, nil
}

// DirEntry fulfills os.DirEntry.
type DirEntry struct {
	EntryName string
	Dir       bool
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

// Sys exists to fulfill os.FileInfo but always returns nil.
func (fi *FileInfo) Sys() any {
	return nil
}

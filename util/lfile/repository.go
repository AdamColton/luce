package lfile

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/upgrade"
)

// File provides an interface fulfilled by *os.File. This allows for testing
// without relying on the actual file system.
type File interface {
	fs.File
	Dir
	io.Writer
}

// Repository is an interface to the file system.
type Repository interface {
	FSOpener
	Create(name string) (fs.File, error)
	Remove(name string) error
	FileStater
	Lstat(name string) (os.FileInfo, error)
	FileReader
}

// OSRepository fulfills Repository by using functions from the "os" package.
type OSRepository struct{}

// Open is a wrapper for os.Open
func (OSRepository) Open(name string) (fs.File, error) {
	f, err := os.Open(name)
	f.Readdirnames(-1)
	return f, err
}

// Create is a wrapper for os.Create
func (OSRepository) Create(name string) (fs.File, error) {
	return os.Create(name)
}

// Remove is a wrapper for os.Remove
func (OSRepository) Remove(name string) error {
	return os.Remove(name)
}

func (OSRepository) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OSRepository) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (OSRepository) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func Size(repo Repository, path string) (int64, error) {
	size := &atomic.Int64{}
	var err error

	// Function to calculate size for a given path
	var calculateSize func(string) *sync.WaitGroup
	calculateSize = func(p string) (wg *sync.WaitGroup) {
		if err != nil {
			return
		}
		fileInfo, localErr := repo.Lstat(p)
		if localErr != nil {
			err = localErr
			return
		}

		// Skip symbolic links to avoid counting them multiple times
		if fileInfo.Mode()&fs.ModeSymlink != 0 {
			return
		}

		if fileInfo.IsDir() {
			dir, localErr := repo.Open(p)
			if localErr != nil {
				err = localErr
				return
			}

			var entries []os.DirEntry
			if drrdr, ok := upgrade.To[DirReader](dir); ok {
				entries, localErr = drrdr.ReadDir(-1)
			}

			dir.Close()
			if localErr != nil {
				err = localErr
				return
			}
			wg = &sync.WaitGroup{}
			wg.Add(len(entries))

			wg = slice.New(entries).Iter().Concurrent(func(entry fs.DirEntry, idx int) {
				innerWg := calculateSize(filepath.Join(p, entry.Name()))
				if innerWg != nil {
					innerWg.Wait()
				}
			})
		} else {
			size.Add(fileInfo.Size())
		}
		return
	}

	// Start calculation from the root path
	wg := calculateSize(path)
	if wg != nil {
		wg.Wait()
	}
	if err != nil {
		return 0, err
	}

	return size.Load(), nil
}

type EmbedFS interface {
	Open(name string) (fs.File, error)
	ReadFile(name string) ([]byte, error)
	ReadDir(name string) ([]fs.DirEntry, error)
}

type embedWrapper struct {
	EmbedFS
}

func (embedWrapper) Create(name string) (fs.File, error) {
	return nil, lerr.Str("embed.FS cannot Create")
}

func (embedWrapper) Lstat(name string) (os.FileInfo, error) {
	return nil, lerr.Str("embed.FS cannot Lstat")
}

func (embedWrapper) Stat(name string) (os.FileInfo, error) {
	return nil, lerr.Str("embed.FS cannot Stat")
}

func (embedWrapper) Remove(name string) error {
	return lerr.Str("embed.FS cannot Remove")
}

type fsFile struct {
	fs.File
	name string
	ew   embedWrapper
}

func (f *fsFile) Name() string {
	return f.name
}

func (f *fsFile) ReadDir(n int) ([]os.DirEntry, error) {
	return f.ew.ReadDir(f.name)
}

func (f *fsFile) Readdirnames(n int) (names []string, err error) {
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

func (f *fsFile) Write(p []byte) (n int, err error) {
	return 0, lerr.Str("embed.FS cannot write")
}

func (ew embedWrapper) Open(name string) (fs.File, error) {
	f, err := ew.EmbedFS.Open(name)
	if err != nil {
		return nil, err
	}
	out := &fsFile{
		File: f,
		name: name,
	}
	return out, err
}

func WrapEmbed(fs EmbedFS) Repository {
	return embedWrapper{fs}
}

type FSReader interface {
	FileStater
	FSOpener
	FileReader
}

type FileStater interface {
	Stat(name string) (os.FileInfo, error)
}

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type FSOpener interface {
	Open(name string) (fs.File, error)
}

func WrapReadFile(fsr FSOpener) func(name string) ([]byte, error) {
	if fr, ok := upgrade.To[FileReader](fsr); ok {
		return fr.ReadFile
	}

	return func(name string) ([]byte, error) {
		f, err := fsr.Open(name)
		if err != nil {
			return nil, err
		}
		return io.ReadAll(f)
	}
}

func WrapStat(fsr FSOpener) func(name string) (os.FileInfo, error) {
	if fr, ok := upgrade.To[FileStater](fsr); ok {
		return fr.Stat
	}

	return func(name string) (os.FileInfo, error) {
		f, err := fsr.Open(name)
		if err != nil {
			return nil, err
		}
		return f.Stat()
	}
}

func WrappedFSReader(o FSOpener) FSReader {
	if r, ok := upgrade.To[FSReader](o); ok {
		return r
	}
	return wrappedFSReader{
		FSOpener: o,
		stat:     WrapStat(o),
		readFile: WrapReadFile(o),
	}
}

type wrappedFSReader struct {
	FSOpener
	readFile func(name string) ([]byte, error)
	stat     func(name string) (os.FileInfo, error)
}

func (wr wrappedFSReader) ReadFile(name string) ([]byte, error) {
	return wr.readFile(name)
}

func (wr wrappedFSReader) Stat(name string) (os.FileInfo, error) {
	return wr.stat(name)
}

func ToReaddirnames(f fs.File) func(n int) (names []string, err error) {
	if rdn, ok := upgrade.To[DirNameReader](f); ok {
		return rdn.Readdirnames
	}

	//TODO: Can I use f.FileInfo.Sys to do this
	return nil
}

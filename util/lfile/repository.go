package lfile

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/adamcolton/luce/ds/slice"
)

// File provides an interface fulfilled by *os.File. This allows for testing
// without relying on the actual file system.
type File interface {
	Dir
	io.Reader
	io.WriteCloser
	Stat() (os.FileInfo, error)
	Readdirnames(n int) (names []string, err error)
	Close() error
}

// Repository is an interface to the file system.
type Repository interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
	Remove(name string) error
	Stat(name string) (os.FileInfo, error)
	Lstat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
}

// OSRepository fulfills Repository by using functions from the "os" package.
type OSRepository struct{}

// Open is a wrapper for os.Open
func (OSRepository) Open(name string) (File, error) {
	return os.Open(name)
}

// Create is a wrapper for os.Create
func (OSRepository) Create(name string) (File, error) {
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

			entries, localErr := dir.ReadDir(-1)
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

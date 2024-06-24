package lfile

import (
	"io"
	"os"
)

// File provides an interface fulfilled by *os.File. This allows for testing
// without relying on the actual file system.
type File interface {
	Dir
	io.Reader
	io.WriteCloser
	Stat() (os.FileInfo, error)
	Readdirnames(n int) (names []string, err error)
}

// Repository is an interface to the file system.
type Repository interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
	Remove(name string) error
	Stat(name string) (os.FileInfo, error)
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

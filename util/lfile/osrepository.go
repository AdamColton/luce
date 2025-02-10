package lfile

import (
	"io/fs"
	"os"
)

// OSRepository fulfills Repository by using functions from the "os" package.
type OSRepository struct{}

// Open is a wrapper for os.Open
func (OSRepository) Open(name string) (fs.File, error) {
	//fmt.Println("Reading ", name, " from OS")
	return os.Open(name)
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

func (OSRepository) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

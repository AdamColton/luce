package lfilemock

import (
	"io/fs"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/adamcolton/luce/ds/lbuf"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/navigator"
)

// Directory is used to create mock directory trees. Directories are
// created when calling Parse.
type Directory struct {
	Name     string
	Children map[string]Node
	// Err is returned by the operations of the File that is opened, and by every
	// operation of a Repository whose root it is.
	Err error
	// Mode is the permission bits Stat reports. Zero means 0755.
	Mode fs.FileMode
	// Mod is the modification time Stat reports.
	Mod time.Time
	// SysData is what Sys returns from the FileInfo.
	SysData any
}

// checkName panics with ErrParseName if name can't be the name of a file or
// directory.
func checkName(name string) {
	if name == "" || name == "." || name == ".." || strings.Contains(name, "/") {
		panic(ErrParseName)
	}
}

// Parse allows for mock directory trees to be setup easily. To create a File
// the value can be either a string or a []byte. To create a sub directory the
// value should be another map[string]any which will be recursivly parsed.
// Adding a []string will create a directory where each file's name and contents
// are the same; if the key is "." the files are added to the directory being
// parsed. A name that is empty, "." (except for that use), ".." or contains a
// "/" panics with ErrParseName, and any other type panics with ErrParseType.
func Parse(root map[string]any) *Directory {
	out := &Directory{
		Children: make(map[string]Node),
	}
	for name, f := range root {
		switch x := f.(type) {
		case map[string]any:
			checkName(name)
			s := Parse(x)
			s.Name = name
			out.Children[name] = s
		case string:
			out.AddFile(name, []byte(x))
		case []string:
			var s *Directory
			if name == "." {
				s = out
			} else {
				checkName(name)
				s = &Directory{
					Name:     name,
					Children: make(map[string]Node, len(x)),
				}
				out.Children[name] = s
			}
			for _, file := range x {
				s.AddFile(file, []byte(file))
			}
		case []byte:
			out.AddFile(name, x)
		default:
			panic(ErrParseType)
		}
	}
	return out
}

// Error fulfills Node and returns Err.
func (d *Directory) Error() error {
	return d.Err
}

// Get finds the Node at a slash separated path below the Directory. Empty and
// "." parts are ignored, so "" is the Directory itself. The Node is the one in
// the tree, so its exported fields, for example Mod, can be changed.
func (d *Directory) Get(name string) (n Node, found bool) {
	n, err := lookup(d, "get", name)
	return n, err == nil
}

// AddDir adds a sub directory. This can be used when modifying the mock
// directory tree. It returns d so calls can be chained. It panics with
// ErrParseName if name is not a valid name.
func (d *Directory) AddDir(name string) *Directory {
	checkName(name)
	c := &Directory{
		Name:     name,
		Children: make(map[string]Node),
	}
	d.Children[name] = c
	return d
}

// AddFile adds a ByteFile. This can be used when modifying the mock directory
// tree. It panics with ErrParseName if name is not a valid name.
func (d *Directory) AddFile(name string, contents []byte) *ByteFile {
	checkName(name)
	f := &ByteFile{
		Name: name,
		Data: lbuf.New(contents),
	}
	d.Children[name] = f
	return f
}

// Repository wraps the Directory in a Repository, which fulfills
// lfile.Repository and treats the Directory as the root.
func (d *Directory) Repository() lfile.Repository {
	return &Repository{
		Node: d,
	}
}

// Next fulfills Node and navigator.Nexter. It returns the child with the name
// key. If there is none and create is true, a Directory is added and returned.
func (d *Directory) Next(key string, create bool, _ navigator.VoidContext) (Node, bool) {
	n, found := d.Children[key]
	if !found && create {
		d.AddDir(key)
		n, found = d.Children[key], true
	}
	return n, found
}

// info describes the directory.
func (d *Directory) info() *FileInfo {
	mode := d.Mode
	if mode == 0 {
		mode = 0o755
	}
	return &FileInfo{
		FileName: d.Name,
		FileMode: mode | fs.ModeDir,
		Mod:      d.Mod,
		Dir:      true,
		SysData:  d.SysData,
	}
}

// File fulfills Node. It creates an instance of *File that fulfills lfile.File.
// The entries are sorted by name.
func (d *Directory) File() *File {
	names := make([]string, 0, len(d.Children))
	for name := range d.Children {
		names = append(names, name)
	}
	sort.Strings(names)
	des := make([]os.DirEntry, len(names))
	for i, name := range names {
		des[i] = d.Children[name].DirEntry()
	}
	info := d.info()
	return &File{
		FileName:   d.Name,
		Dir:        true,
		DirEntries: des,
		FileInfo:   info,
		FileMode:   info.FileMode,
		Mod:        d.Mod,
		SysData:    d.SysData,
		Err:        d.Err,
	}
}

// DirEntry fulfills Node. It creates a DirEntry instance for the Directory.
func (d *Directory) DirEntry() os.DirEntry {
	info := d.info()
	return &DirEntry{
		EntryName: d.Name,
		Dir:       true,
		FileMode:  info.FileMode.Type(),
		FileInfo:  info,
	}
}

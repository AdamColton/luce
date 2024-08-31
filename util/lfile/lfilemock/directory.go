package lfilemock

import (
	"os"
	"strings"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/navigator"
)

// Directory is used to create mock directory trees. Directories are
// created when calling ParseDir.
type Directory struct {
	Name     string
	Children map[string]Node
	Err      error
}

// ParseDir allows for mock directory trees to be setup easily. To create a File
// the value can be either a string or a []byte. To create a sub directory the
// value should be another map[string]any which will be recursivly parsed.
// Adding a []string will create a directory where each file's name and contents
// are the same.
func Parse(root map[string]any) *Directory {
	// TODO: handle bad characters like /
	out := &Directory{
		Children: make(map[string]Node),
	}
	for name, f := range root {
		switch x := f.(type) {
		case map[string]any:
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
func (d *Directory) Error() error {
	return d.Err
}

func (d *Directory) Get(path string) (n Node, found bool) {
	vc := navigator.VoidContext{}
	n = d
	for _, key := range strings.Split(path, "/") {
		n, found = n.Next(key, false, vc)
		if !found {
			break
		}
	}
	if found {
		n = pathNode{
			name: path,
			Node: n,
		}
	}
	return
}

// AddDir adds a sub directory. This can be used when modifying the mock
// directory tree.
func (d *Directory) AddDir(name string) *Directory {
	c := &Directory{
		Name:     name,
		Children: make(map[string]Node),
	}
	d.Children[name] = c
	return d
}

// AddFile adds a ByteFile. This can be used when modifying the mock directory
// tree.
func (d *Directory) AddFile(name string, contents []byte) *ByteFile {
	f := &ByteFile{
		Name: name,
		Data: contents,
	}
	d.Children[name] = f
	return f
}

// Parse allows for mock directory trees to be setup easily. To create a
// File the value can be either a string or a []byte. To create a sub directory
// the value should be another map[string]any which will be recursivly parsed.
// The returned Repository fulfills lfile.Repository.
func (d *Directory) Repository() lfile.Repository {
	return &Repository{
		Node: d,
	}
}

// Next fulfills Node and navigator.Nexter. It is used internally for navigating
// the mock directory.
func (d *Directory) Next(key string, create bool, _ navigator.VoidContext) (Node, bool) {
	n, found := d.Children[key]
	if !found && create {
		n = d.AddDir(key)
		found = true
	}
	return n, found
}

// File fulfills Node. It creates an instance of *File that fulfills lfile.File.
func (d *Directory) File() *File {
	des := make([]os.DirEntry, 0, len(d.Children))
	for _, c := range d.Children {
		des = append(des, c.DirEntry())
	}
	return &File{
		FileName:   d.Name,
		Dir:        true,
		DirEntries: des,
		FileInfo: &FileInfo{
			FileName: d.Name,
			Dir:      true,
		},
		FileMode: 0,
		Err:      d.Err,
	}
}

// DirEntry fulfills Node. It creates a DirEntry instance for the Directory.
func (d *Directory) DirEntry() os.DirEntry {
	return &DirEntry{
		EntryName: d.Name,
		Dir:       true,
		FileInfo: &FileInfo{
			FileName: d.Name,
			Dir:      true,
		},
	}
}

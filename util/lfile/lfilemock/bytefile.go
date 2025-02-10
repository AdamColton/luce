package lfilemock

import (
	"io/fs"
	"os"
	"time"

	"github.com/adamcolton/luce/ds/lbuf"
	"github.com/adamcolton/luce/util/navigator"
)

// ByteFile is used to create files in mock directory trees. ByteFiles are
// created when calling Parse.
type ByteFile struct {
	Name string
	Data *lbuf.Buffer
	// Err is returned by the operations of the File that is opened.
	Err error
	// Mode is the permission bits Stat reports. Zero means 0644.
	Mode fs.FileMode
	// Mod is the modification time Stat reports.
	Mod time.Time
	// SysData is what Sys returns from the FileInfo, for example a
	// *syscall.Stat_t.
	SysData any
}

// info describes the file.
func (f *ByteFile) info() *FileInfo {
	mode := f.Mode
	if mode == 0 {
		mode = 0o644
	}
	return &FileInfo{
		FileName: f.Name,
		FileSize: int64(f.Data.Len()),
		FileMode: mode &^ fs.ModeDir,
		Mod:      f.Mod,
		SysData:  f.SysData,
	}
}

// File fulfills Node. It uses a ByteFile to create an instance of File that
// fulfills lfile.File. The File reads from the start of the data, and what is
// written to it is written to the ByteFile.
func (f *ByteFile) File() *File {
	f.Data.Idx = 0
	info := f.info()
	return &File{
		FileName: f.Name,
		Buffer:   f.Data,
		Dir:      false,
		FileSize: info.FileSize,
		FileInfo: info,
		FileMode: info.FileMode,
		Mod:      f.Mod,
		SysData:  f.SysData,
		Err:      f.Err,
	}
}

// DirEntry fulfills Node. It creates a DirEntry instance for the ByteFile.
func (f *ByteFile) DirEntry() os.DirEntry {
	return &DirEntry{
		EntryName: f.Name,
		Dir:       false,
		FileInfo:  f.info(),
	}
}

// Next fulfills Node and navigator.Nexter. A file has no children, so it always
// returns false.
func (f *ByteFile) Next(key string, create bool, _ navigator.VoidContext) (Node, bool) {
	return nil, false
}

// Error fulfills Node and returns Err.
func (f *ByteFile) Error() error {
	return f.Err
}

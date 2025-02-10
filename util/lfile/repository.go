package lfile

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/upgrade"
)

func Size(o FSOpener, path string) (int64, error) {
	size := &atomic.Int64{}
	var err error
	lstat := WrapLstat(o)

	// Function to calculate size for a given path
	var calculateSize func(string) *sync.WaitGroup
	calculateSize = func(p string) (wg *sync.WaitGroup) {
		if err != nil {
			return
		}
		fileInfo, localErr := lstat(p)
		if localErr != nil {
			err = localErr
			return
		}

		// Skip symbolic links to avoid counting them multiple times
		if fileInfo.Mode()&fs.ModeSymlink != 0 {
			return
		}

		if fileInfo.IsDir() {
			dir, localErr := o.Open(p)
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

func WrapLstat(fsr FSOpener) func(name string) (os.FileInfo, error) {
	if fr, ok := upgrade.To[FSLstater](fsr); ok {
		return fr.Lstat
	}

	return func(name string) (os.FileInfo, error) {
		f, err := fsr.Open(name)
		if err != nil {
			return nil, err
		}
		if ls, ok := upgrade.To[FileLstater](f); ok {
			return ls.Lstat()
		}
		return f.Stat()
	}
}

// func WrapReadFile(fsr FSOpener) func(name string) ([]byte, error) {
// 	if fr, ok := upgrade.To[FileReader](fsr); ok {
// 		return fr.ReadFile
// 	}

// 	return func(name string) ([]byte, error) {
// 		f, err := fsr.Open(name)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return io.ReadAll(f)
// 	}
// }

// func WrapStat(fsr FSOpener) func(name string) (os.FileInfo, error) {
// 	if fr, ok := upgrade.To[FSStater](fsr); ok {
// 		return fr.Stat
// 	}

// 	return func(name string) (os.FileInfo, error) {
// 		f, err := fsr.Open(name)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return f.Stat()
// 	}
// }

// func WrappedFSReader(o FSOpener) FSReader {
// 	if r, ok := upgrade.To[FSReader](o); ok {
// 		return r
// 	}
// 	return wrappedFSReader{
// 		FSOpener: o,
// 		stat:     WrapStat(o),
// 		readFile: WrapReadFile(o),
// 	}
// }

// type wrappedFSReader struct {
// 	FSOpener
// 	readFile func(name string) ([]byte, error)
// 	stat     func(name string) (os.FileInfo, error)
// 	lstat    func(name string) (os.FileInfo, error)
// }

// func (wr wrappedFSReader) ReadFile(name string) ([]byte, error) {
// return wr.readFile(name)
// }
//
// func (wr wrappedFSReader) Stat(name string) (os.FileInfo, error) {
// return wr.stat(name)
// }
//
// func (wr wrappedFSReader) Lstat(name string) (os.FileInfo, error) {
// return wr.lstat(name)
// }

//func ToReaddirnames(f fs.File) func(n int) (names []string, err error) {
//	if rdn, ok := upgrade.To[DirNameReader](f); ok {
//		return rdn.Readdirnames
//	}
//
//	//TODO: Can I use f.FileInfo.Sys to do this
//	return nil
//}

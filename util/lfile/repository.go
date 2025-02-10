package lfile

import (
	"io/fs"
	"path"
	"sync"
	"sync/atomic"

	"github.com/adamcolton/luce/util/upgrade"
)

// Size returns the size in bytes of the file at root. For a directory it is the
// sum of the sizes of everything under it. Symbolic links are skipped, so
// nothing is counted twice; a file system that can't Lstat (see WrapLstat) can't
// tell a link from what it points to. Directories are read concurrently. If more
// than one read fails, one of the errors is returned, and which one is not
// defined.
func Size(o FSOpener, root string) (int64, error) {
	var (
		size     atomic.Int64
		mux      sync.Mutex
		firstErr error
	)
	fail := func(err error) {
		mux.Lock()
		if firstErr == nil {
			firstErr = err
		}
		mux.Unlock()
	}
	failed := func() bool {
		mux.Lock()
		defer mux.Unlock()
		return firstErr != nil
	}
	lstat := WrapLstat(o)

	var walk func(p string)
	walk = func(p string) {
		if failed() {
			return
		}
		info, err := lstat(p)
		if err != nil {
			fail(err)
			return
		}

		if info.Mode()&fs.ModeSymlink != 0 {
			return
		}
		if !info.IsDir() {
			size.Add(info.Size())
			return
		}

		entries, err := fs.ReadDir(o, p)
		if err != nil {
			fail(err)
			return
		}
		var wg sync.WaitGroup
		for _, entry := range entries {
			child := path.Join(p, entry.Name())
			if !entry.IsDir() {
				walk(child)
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				walk(child)
			}()
		}
		wg.Wait()
	}

	walk(root)
	if firstErr != nil {
		return 0, firstErr
	}
	return size.Load(), nil
}

// WrapLstat returns a function that describes a file without following a
// symbolic link. It is the Lstat of the file system if it has one. Otherwise it
// opens the file and uses its Lstat if it has one; failing that it uses Stat,
// which follows links, and so a link can't be told from its target.
func WrapLstat(fsr FSOpener) func(name string) (fs.FileInfo, error) {
	if fr, ok := upgrade.To[FSLstater](fsr); ok {
		return fr.Lstat
	}

	return func(name string) (fs.FileInfo, error) {
		f, err := fsr.Open(name)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if ls, ok := upgrade.To[FileLstater](f); ok {
			return ls.Lstat()
		}
		return f.Stat()
	}
}

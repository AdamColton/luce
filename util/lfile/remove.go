package lfile

import (
	"errors"
	"io/fs"
)

// TryRemove removes name, and it is not an error if there is nothing there to
// remove. Any other error, for example a directory that is not empty, is
// returned. It is for clearing a path before using it, as a socket does with the
// file a previous run left behind.
func TryRemove(r FSRemover, name string) error {
	if err := r.Remove(name); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

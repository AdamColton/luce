package lfile

// IterHandler represents something that will handle each value in the iterator.
type IterHandler interface {
	HandleIter(Iterator)
}

// RunHandler will create an Iter from Iterator and call the HandleIter method
// on the IterHandler for each value in the iterator.
func RunHandlerSource(ii IteratorSource, ih IterHandler) error {
	i, done := ii.Iterator()
	for ; !done; _, done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}

// RunHandler will create an Iter from Iterator and call the HandleIter method
// on the IterHandler for each value in the iterator.
func RunHandler(i Iterator, ih IterHandler) error {
	for done := i.Reset(); !done; _, done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}

// GetByTypeHandler records all the files and directories the Iterator visits
// and seperates them by type.
type GetByTypeHandler struct {
	Files, Dirs []string
}

// HandleIter fulfills IterHandler and records the current location based on
// the type.
func (bt *GetByTypeHandler) HandleIter(i Iterator) {
	if i.Stat().IsDir() {
		bt.Dirs = append(bt.Dirs, i.Path())
	} else {
		bt.Files = append(bt.Files, i.Path())
	}
}

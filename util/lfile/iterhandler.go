package lfile

// IterHandler represents something that will handle each value in the iterator.
type IterHandler interface {
	HandleIter(Iterator)
}

// RunHandlerSource gets an Iterator from ii and calls HandleIter on ih for each
// value in it. It returns the error that ended the iteration, if there was one.
func RunHandlerSource(ii IteratorSource, ih IterHandler) error {
	i, done := ii.Iterator()
	for ; !done; _, done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}

// RunHandler resets i and calls HandleIter on ih for each value in it. It
// returns the error that ended the iteration, if there was one.
func RunHandler(i Iterator, ih IterHandler) error {
	for done := i.Reset(); !done; _, done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}

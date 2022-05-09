package lfile

// IterHandler represents something that will handle each value in the iterator.
type IterHandler interface {
	HandleIter(Iterator)
}

// RunHandler will creat an Iter from Iterator and call the HandleIter method on
// the IterHandler for each value in the iterator.
func RunHandlerSource(ii IteratorSource, ih IterHandler) error {
	i, done := ii.Iterator()
	for ; !done; done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}

// RunHandler will creat an Iter from Iterator and call the HandleIter method on
// the IterHandler for each value in the iterator.
func RunHandler(i Iterator, ih IterHandler) error {
	for done := i.Reset(); !done; done = i.Next() {
		ih.HandleIter(i)
	}
	return i.Err()
}
